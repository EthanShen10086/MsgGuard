package memory

import "context"

type Queue struct{}

func NewQueue() *Queue { return &Queue{} }

func (q *Queue) Publish(ctx context.Context, subject string, data []byte) error {
	return nil
}
