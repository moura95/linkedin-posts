package main

import (
	"errors"
	"fmt"

	strategy "linkedinPosts/implement_factory_and_strategy/strategies"
)

var (
	NotImplementedErr     = errors.New("not implemented")
	UnsupportedChannelErr = errors.New("unsupported notification channel")
)

type NotificationFactory interface {
	CreateStrategy(channelType string) (strategy.NotificationStrategy, error)
}
type notificationFactory struct{}

func (f *notificationFactory) CreateStrategy(channelType string) (strategy.NotificationStrategy, error) {
	switch channelType {
	case "email":
		return &strategy.EmailStrategy{}, nil
	case "sms":
		return &strategy.SMSStrategy{}, nil
	case "push":
		return &strategy.PushStrategy{}, nil
	case "slack":
		return &strategy.SlackStrategy{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", UnsupportedChannelErr, channelType)
	}
}
