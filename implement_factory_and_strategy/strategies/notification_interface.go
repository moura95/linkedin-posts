package strategy

import "errors"

var (
	NotImplementedErr     = errors.New("not implemented")
	UnsupportedChannelErr = errors.New("unsupported notification channel")
)

type NotificationStrategy interface {
	Send(message string, recipient string) error
	GetChannelName() string
}

type BaseNotificationStrategy struct{}

func (s *BaseNotificationStrategy) Send(message string, recipient string) error {
	return NotImplementedErr
}

func (s *BaseNotificationStrategy) GetChannelName() string {
	return "Base"
}
