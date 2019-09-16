package apns

import (
	"context"
	"fmt"
)

// NewMock returns mock of `Notificator`.
func NewMock() Notificator {
	return &mock{}
}

type mock struct{}

func (m *mock) Push(ctx context.Context, token string) error {
	fmt.Println("sending update notification to", token)
	return nil
}
