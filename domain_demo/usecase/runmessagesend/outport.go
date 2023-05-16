package runmessagesend

import "context"

type Outport interface {
	SendMessage(ctx context.Context, message string) (string, error)
}
