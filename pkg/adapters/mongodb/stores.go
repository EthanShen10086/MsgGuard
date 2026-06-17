package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

const defaultDB = "msgguard"

func connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return client, nil
}

type FeedbackStore struct {
	col *mongo.Collection
}

func NewFeedbackStore(uri string) (*FeedbackStore, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}
	col := client.Database(defaultDB).Collection("feedback")
	return &FeedbackStore{col: col}, nil
}

func (s *FeedbackStore) Create(ctx context.Context, item ports.FeedbackItem) error {
	_, err := s.col.InsertOne(ctx, bson.M{
		"id": item.ID, "body": item.Body, "label": item.Label,
		"locale": item.Locale, "tenant_id": item.TenantID, "trace_id": item.TraceID, "created_at": item.CreatedAt,
	})
	return err
}

func (s *FeedbackStore) List(ctx context.Context, limit int) ([]ports.FeedbackItem, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cur, err := s.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []ports.FeedbackItem
	for cur.Next(ctx) {
		var doc struct {
			ID        string    `bson:"id"`
			Body      string    `bson:"body"`
			Label     string    `bson:"label"`
			Locale    string    `bson:"locale"`
			TenantID  string    `bson:"tenant_id"`
			TraceID   string    `bson:"trace_id"`
			CreatedAt time.Time `bson:"created_at"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, ports.FeedbackItem{
			ID: doc.ID, Body: doc.Body, Label: doc.Label,
			Locale: doc.Locale, TenantID: doc.TenantID, TraceID: doc.TraceID, CreatedAt: doc.CreatedAt,
		})
	}
	return out, cur.Err()
}

type RuleStore struct {
	col *mongo.Collection
}

func NewRuleStore(uri string) (*RuleStore, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}
	col := client.Database(defaultDB).Collection("rules")
	return &RuleStore{col: col}, nil
}

func (s *RuleStore) GetLatest(ctx context.Context) (*ports.RuleBundle, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var doc struct {
		Version  string    `bson:"version"`
		Checksum string    `bson:"checksum"`
		Payload  []byte    `bson:"payload"`
		Created  time.Time `bson:"created_at"`
	}
	err := s.col.FindOne(ctx, bson.M{}, opts).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return &ports.RuleBundle{Version: "1.0.0", Checksum: "seed", Payload: []byte(`{"keywords":["免费"]}`)}, nil
	}
	if err != nil {
		return nil, err
	}
	return &ports.RuleBundle{Version: doc.Version, Checksum: doc.Checksum, Payload: doc.Payload}, nil
}

func (s *RuleStore) GetByVersion(ctx context.Context, version string) (*ports.RuleBundle, error) {
	var doc struct {
		Version  string `bson:"version"`
		Checksum string `bson:"checksum"`
		Payload  []byte `bson:"payload"`
	}
	err := s.col.FindOne(ctx, bson.M{"version": version}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("version not found")
	}
	if err != nil {
		return nil, err
	}
	return &ports.RuleBundle{Version: doc.Version, Checksum: doc.Checksum, Payload: doc.Payload}, nil
}

func (s *RuleStore) Save(ctx context.Context, bundle ports.RuleBundle) error {
	if !json.Valid(bundle.Payload) {
		bundle.Payload = []byte(`{}`)
	}
	opts := options.Update().SetUpsert(true)
	_, err := s.col.UpdateOne(ctx,
		bson.M{"version": bundle.Version},
		bson.M{"$set": bson.M{
			"version": bundle.Version, "checksum": bundle.Checksum,
			"payload": bundle.Payload, "created_at": time.Now().UTC(),
		}},
		opts,
	)
	return err
}

type AnalyticsStore struct {
	col *mongo.Collection
}

func NewAnalyticsStore(uri string) (*AnalyticsStore, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}
	col := client.Database(defaultDB).Collection("analytics_events")
	return &AnalyticsStore{col: col}, nil
}

func (s *AnalyticsStore) Insert(ctx context.Context, event ports.AnalyticsEvent) error {
	ts := event.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := s.col.InsertOne(ctx, bson.M{
		"id": event.ID, "name": event.Name, "props": event.Props,
		"device_id": event.DeviceID, "tenant_id": event.TenantID,
		"trace_id": event.TraceID, "created_at": ts,
	})
	return err
}

func (s *AnalyticsStore) List(ctx context.Context, since time.Time, limit int) ([]ports.AnalyticsEvent, error) {
	filter := bson.M{}
	if !since.IsZero() {
		filter["created_at"] = bson.M{"$gte": since}
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	cur, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	return scanAnalytics(cur, ctx)
}

func (s *AnalyticsStore) CountByName(ctx context.Context, since time.Time) (map[string]int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": since}}}},
		{{Key: "$group", Value: bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}}},
	}
	cur, err := s.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	counts := map[string]int{}
	for cur.Next(ctx) {
		var doc struct {
			Name  string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		counts[doc.Name] = doc.Count
	}
	return counts, cur.Err()
}

func (s *AnalyticsStore) DeleteByDeviceID(ctx context.Context, deviceID string) (int, error) {
	res, err := s.col.DeleteMany(ctx, bson.M{"device_id": deviceID})
	if err != nil {
		return 0, err
	}
	return int(res.DeletedCount), nil
}

func scanAnalytics(cur *mongo.Cursor, ctx context.Context) ([]ports.AnalyticsEvent, error) {
	var out []ports.AnalyticsEvent
	for cur.Next(ctx) {
		var doc struct {
			ID        string         `bson:"id"`
			Name      string         `bson:"name"`
			Props     map[string]any `bson:"props"`
			DeviceID  string         `bson:"device_id"`
			TenantID  string         `bson:"tenant_id"`
			TraceID   string         `bson:"trace_id"`
			CreatedAt time.Time      `bson:"created_at"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, ports.AnalyticsEvent{
			ID: doc.ID, Name: doc.Name, Props: doc.Props,
			DeviceID: doc.DeviceID, TenantID: doc.TenantID,
			TraceID: doc.TraceID, Timestamp: doc.CreatedAt,
		})
	}
	return out, cur.Err()
}
