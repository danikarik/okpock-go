package memory_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/mail"
	"github.com/danikarik/okpock/pkg/mail/memory"
	"github.com/stretchr/testify/assert"
)

var requiredVars = []string{
	"TEST_RECIPIENT",
}

func TestSendMail(t *testing.T) {
	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	mailer := memory.New()

	testCase := struct {
		HTML string
		Text string
	}{
		HTML: "<h1>Mock Test Email</h1>",
		Text: "This email was sent with Mock Mailer.",
	}

	message := mail.NewMessage(
		"noreply@okpock.com",
		env.Get("TEST_RECIPIENT"),
		"Test Message From Unit Test",
		testCase.HTML,
		testCase.Text,
		mail.DefaultCharset,
	)

	sentAt, err := mailer.SendMail(ctx, message)
	assert.NoError(err)
	assert.NotEmpty(message.ID)
	assert.NotNil(sentAt)
}
