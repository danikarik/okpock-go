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
<p><a href="{{ .ConfirmationURL }}">Confirm your email</a></p>
`

const confirmText = `
Confirm your signup

Follow this link to confirm your account:
{{ .ConfirmationURL }}
`

const inviteHTML = `
<h2>Accept invite</h2>

<p>{{ .Referer }} has invited you to join {{ .AppURL }} workspace.</p>
<p>Follow this link to confirm your account:</p>
<p><a href="{{ .ConfirmationURL }}">Confirm your email</a></p>
`

const inviteText = `
Accept invite

{{ .Referer }} has invited you to join {{ .AppURL }} workspace.
Follow this link to confirm your account:
{{ .ConfirmationURL }}
`

const emailChangeHTML = `
<h2>Confirm your email change</h2>

<p>Follow this link to confirm your account:</p>
<p><a href="{{ .ConfirmationURL }}">Confirm your email</a></p>
`

const emailChangeText = `
Confirm your email change

Follow this link to confirm your email change:
{{ .ConfirmationURL }}
`

func (s *Service) mailBody(htmlName, htmlContent, textName, textContent string,
	data interface{}) (string, string, error) {

	var (
		htmlBody bytes.Buffer
		textBody bytes.Buffer
	)

	html, err := htmlTemplate.New(htmlName).Parse(htmlContent)
	if err != nil {
		return "", "", err
	}

	err = html.Execute(&htmlBody, data)
	if err != nil {
		return "", "", err
	}

	text, err := textTemplate.New(textName).Parse(textContent)
	if err != nil {
		return "", "", err
	}

	err = text.Execute(&textBody, data)
	if err != nil {
		return "", "", err
	}

	return htmlBody.String(), textBody.String(), nil
}

func (s *Service) recoverMessage(u *api.User) (*mail.Message, error) {
	url, err := s.confirmationURL(u, api.RecoveryConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{"ConfirmationURL": url}

	html, text, err := s.mailBody(
		"recover_html",
		recoverHTML,
		"recover_text",
		recoverText,
		data,
	)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.Email,
		"Reset Password",
		html,
		text,
		mail.DefaultCharset,
	), nil
}

func (s *Service) confirmMessage(u *api.User) (*mail.Message, error) {
	url, err := s.confirmationURL(u, api.SignUpConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{"ConfirmationURL": url}

	html, text, err := s.mailBody(
		"confirm_html",
		confirmHTML,
		"confirm_text",
		confirmText,
		data,
	)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.Email,
		"Confirm signup",
		html,
		text,
		mail.DefaultCharset,
	), nil
}

func (s *Service) inviteMessage(ref, u *api.User) (*mail.Message, error) {
	url, err := s.confirmationURL(u, api.InviteConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{
		"ConfirmationURL": url,
		"Referer":         ref.Email,
		"AppURL":          s.appURL(""),
	}

	html, text, err := s.mailBody(
		"invite_html",
		inviteHTML,
		"invite_text",
		inviteText,
		data,
	)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.Email,
		"Accept invite",
		html,
		text,
		mail.DefaultCharset,
	), nil
}

func (s *Service) emailChangeMessage(u *api.User) (*mail.Message, error) {
	url, err := s.confirmationURL(u, api.EmailChangeConfirmation)
	if err != nil {
		return nil, err
	}

	data := M{"ConfirmationURL": url}

	html, text, err := s.mailBody(
		"email_change_html",
		emailChangeHTML,
		"email_change_text",
		emailChangeText,
		data,
	)
	if err != nil {
		return nil, err
	}

	return mail.NewMessage(
		defaultSender,
		u.EmailChange,
		"Email change",
		html,
		text,
		mail.DefaultCharset,
	), nil
}
