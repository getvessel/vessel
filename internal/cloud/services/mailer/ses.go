package mailer

type SESMailer struct{}

func NewSESMailer() *SESMailer {
	return &SESMailer{}
}

func (m *SESMailer) SendEmail(to, subject, body string) error {
	// Implementation for AWS SES
	return nil
}
