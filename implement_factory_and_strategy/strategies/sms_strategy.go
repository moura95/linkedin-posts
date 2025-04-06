package strategy

import (
	"fmt"
	"time"
)

type SMSStrategy struct {
	BaseNotificationStrategy
}

func (s *SMSStrategy) Send(message string, recipient string) error {
	fmt.Printf("SMS: Sending to %s: %s\n", recipient, message)
	time.Sleep(200 * time.Millisecond)
	return nil
}

func (s *SMSStrategy) GetChannelName() string {
	return "SMS"
}
