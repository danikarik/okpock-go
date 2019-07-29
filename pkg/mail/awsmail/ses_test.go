package awsmail_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/mail"
	"github.com/danikarik/okpock/pkg/mail/awsmail"
	"github.com/stretchr/testify/assert"
)

var requiredVars = []string{
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"MAILER_REGION",
	"TEST_RECIPIENT",
}

func skipTest(t *testing.T) {
	if v, ok := os.LookupEnv("SKIP_SES_TEST"); ok {
		skip, err := strconv.ParseBool(v)
		if err == nil && skip {
			t.Skip(`skip test: SKIP_SES_TEST is present`)
		}
	}
}

func TestSendMail(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	mailer, err := awsmail.New(env.Get("MAILER_REGION"))
	if !assert.NoError(err) {
		return
	}

	testCase := struct {
		HTML string
		Text string
	}{
		HTML: "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>",
		Text: "This email was sent with Amazon SES using the AWS SDK for Go.",
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
