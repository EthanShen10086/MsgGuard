package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type SubscriptionStore struct {
	col *mongo.Collection
}

func NewSubscriptionStore(uri string) (*SubscriptionStore, error) {
	client, err := connect(uri)
	if err != nil {
		return nil, err
	}
	col := client.Database(defaultDB).Collection("subscriptions")
	return &SubscriptionStore{col: col}, nil
}

func (s *SubscriptionStore) Upsert(ctx context.Context, sub ports.Subscription) error {
	ts := sub.UpdatedAt
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	doc := bson.M{
		"device_id": sub.DeviceID, "product_id": sub.ProductID,
		"signed_transaction": sub.SignedTransaction, "is_pro": sub.IsPro,
		"expires_at": sub.ExpiresAt, "updated_at": ts,
	}
	opts := options.Update().SetUpsert(true)
	_, err := s.col.UpdateOne(ctx, bson.M{"device_id": sub.DeviceID}, bson.M{"$set": doc}, opts)
	return err
}

func (s *SubscriptionStore) GetByDeviceID(ctx context.Context, deviceID string) (*ports.Subscription, error) {
	var doc struct {
		DeviceID          string     `bson:"device_id"`
		ProductID         string     `bson:"product_id"`
		SignedTransaction string     `bson:"signed_transaction"`
		IsPro             bool       `bson:"is_pro"`
		ExpiresAt         *time.Time `bson:"expires_at"`
		UpdatedAt         time.Time  `bson:"updated_at"`
	}
	err := s.col.FindOne(ctx, bson.M{"device_id": deviceID}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return &ports.Subscription{DeviceID: deviceID, IsPro: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return &ports.Subscription{
		DeviceID: doc.DeviceID, ProductID: doc.ProductID,
		SignedTransaction: doc.SignedTransaction, IsPro: doc.IsPro,
		ExpiresAt: doc.ExpiresAt, UpdatedAt: doc.UpdatedAt,
	}, nil
}
