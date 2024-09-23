package smsservice

import (
	"log/slog"
	"testing"
)

func Test_sendSMS(t *testing.T) {
	// Create a new Twilio SMS sender
	config := &TwilioSMSSenderConfig{
		AccountSID: "ACd08507e1b3f7f31a2e72932f10d3564d",
		AuthToken:  "52f59f5eee8e6cd0a8e05fddd6d450df",
		FromNumber: "+14042366595",
	}

	sender := NewTwilioSMSSender(config, &slog.Logger{})

	// Send a SMS
	err := sender.SendSMS("+4915170640522", "Hello, World!")
	if err != nil {
		t.Errorf("Error sending SMS: %v", err)
	}

	t.FailNow()
}
