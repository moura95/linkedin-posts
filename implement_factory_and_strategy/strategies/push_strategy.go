package strategy

import (
	"fmt"
	"time"
)

type PushStrategy struct {
	BaseNotificationStrategy
}

func (s *PushStrategy) Send(message string, recipient string) error {
	fmt.Printf("PUSH: Send to device %s: %s\n", recipient, message)
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (s *PushStrategy) GetChannelName() string {
	return "Push"
}
