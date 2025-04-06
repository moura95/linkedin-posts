package strategy

import (
	"fmt"
	"time"
)

type SlackStrategy struct {
	BaseNotificationStrategy
}

func (s *SlackStrategy) Send(message string, recipient string) error {
	fmt.Printf("SLACK: Send to channel %s: %s\n", recipient, message)
	time.Sleep(150 * time.Millisecond)
	return nil
}

func (s *SlackStrategy) GetChannelName() string {
	return "Slack"
}
