package strategy

import (
	"fmt"
	"time"
)

type EmailStrategy struct {
	BaseNotificationStrategy
}

func (s *EmailStrategy) Send(message string, recipient string) error {
	fmt.Printf("EMAIL: Sending to %s: %s\n", recipient, message)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *EmailStrategy) GetChannelName() string {
	return "Email"
}
