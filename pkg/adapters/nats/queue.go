package natsadapter

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Queue struct {
	conn *nats.Conn
}

func NewQueue(url string) (*Queue, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Queue{conn: nc}, nil
}

func (q *Queue) Publish(ctx context.Context, subject string, data []byte) error {
	return q.conn.Publish(subject, data)
}

func (q *Queue) Close() error {
	if q.conn != nil {
		q.conn.Close()
	}
	return nil
}
