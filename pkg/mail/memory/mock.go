package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danikarik/okpock/pkg/mail"
	uuid "github.com/satori/go.uuid"
)

// New returns mock mailer.
func New() mail.Mailer {
	return &mockMailer{}
}

type mockMailer struct{}

func (m *mockMailer) SendMail(ctx context.Context, message *mail.Message) (*time.Time, error) {
	message.ID = uuid.NewV4().String()
	now := time.Now()

	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))

	return &now, nil
}
