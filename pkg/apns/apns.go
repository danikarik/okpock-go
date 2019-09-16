package apns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

// New creates a new instance of `Notificator`.
func New(data []byte, pass string, production bool) (Notificator, error) {
	cert, err := certificate.FromP12Bytes(data, pass)
	if err != nil {
		return nil, err
	}
	client := apns2.NewClient(cert)
	if production {
		client = client.Production()
	}
	return &notificator{client}, nil
}

// Notificator sends notifictions for iOS devices.
type Notificator interface {
	Push(ctx context.Context, token string) error
}

type notificator struct {
	client *apns2.Client
}

func (n *notificator) Push(ctx context.Context, token string) error {
	result, err := n.client.PushWithContext(ctx, &apns2.Notification{
		DeviceToken: token,
		Payload:     []byte(`{}`),
		Priority:    apns2.PriorityHigh,
	})
	if err != nil {
		return err
	}
	if result.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code: %s", result.Reason)
	}
	return nil
}
