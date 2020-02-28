package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	m := NewMailer()
	err := m.SendEmail(EmailRegistration, struct {
		Name string
		URL  string
	}{
		"Akexander",
		"http://ya.ru",
	}, "bopohuh@ya.ru")
	assert.NoError(t, err)
}
