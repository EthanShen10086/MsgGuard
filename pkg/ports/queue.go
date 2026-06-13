package ports

import "context"

type Queue interface {
	Publish(ctx context.Context, subject string, data []byte) error
}
