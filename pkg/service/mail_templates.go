package service

import (
	"bytes"
	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/mail"
)

const defaultSender = "noreply@okpock.com"

const recoverHTML = `
<h2>Reset Password</h2>

<p>Follow this link to reset the password for your account:</p>
<p><a href="{{ .ConfirmationURL }}">Reset Password</a></p>
`

const recoverText = `
Reset Password

Follow this link to reset the password for your account:
{{ .ConfirmationURL }}
`

const confirmHTML = `
<h2>Confirm your signup</h2>

<p>Follow this link to confirm your account:</p>
<p><a href="{{ .ConfirmationURL }}">Confirm your mail</a></p>
`

const confirmText = `
Confirm your signup

Follow this link to confirm your account:
{{ .ConfirmationURL }}
`

func (s *Service) recoverMessage(u *api.User) (*mail.Message, error) {
	var (
		htmlBody bytes.Buffer
		textBody bytes.Buffer
	)

	url, err := s.confirmationURL(u, api.RecoveryConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{"ConfirmationURL": url}

	html, err := htmlTemplate.New("recover_html").Parse(recoverHTML)
	if err != nil {
		return nil, err
	}

	err = html.Execute(&htmlBody, data)
	if err != nil {
		return nil, err
	}

	text, err := textTemplate.New("recover_text").Parse(recoverText)
	if err != nil {
		return nil, err
	}

	err = text.Execute(&textBody, data)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.Email,
		"Reset Password",
		htmlBody.String(),
		textBody.String(),
		mail.DefaultCharset,
	), nil
}

func (s *Service) confirmMessage(u *api.User) (*mail.Message, error) {
	var (
		htmlBody bytes.Buffer
		textBody bytes.Buffer
	)

	url, err := s.confirmationURL(u, api.RecoveryConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{"ConfirmationURL": url}

	html, err := htmlTemplate.New("confirm_html").Parse(confirmHTML)
	if err != nil {
		return nil, err
	}

	err = html.Execute(&htmlBody, data)
	if err != nil {
		return nil, err
	}

	text, err := textTemplate.New("confirm_text").Parse(confirmText)
	if err != nil {
		return nil, err
	}

	err = text.Execute(&textBody, data)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.Email,
		"Confirm signup",
		htmlBody.String(),
		textBody.String(),
		mail.DefaultCharset,
	), nil
}
