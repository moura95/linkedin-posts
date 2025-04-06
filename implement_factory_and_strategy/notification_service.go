package main

import strategy "linkedinPosts/implement_factory_and_strategy/strategies"

type NotificationService interface {
	SetStrategy(channelType string) error
	Notify(message string, recipient string) error
	GetCurrentChannel() string
}

type DefaultNotificationService struct {
	strategy strategy.NotificationStrategy
	factory  NotificationFactory
}

func NewNotificationService(factory NotificationFactory, defaultChannel string) (NotificationService, error) {
	strategy, err := factory.CreateStrategy(defaultChannel)
	if err != nil {
		return nil, err
	}

	return &DefaultNotificationService{
		strategy: strategy,
		factory:  factory,
	}, nil
}

func (s *DefaultNotificationService) SetStrategy(channelType string) error {
	strategy, err := s.factory.CreateStrategy(channelType)
	if err != nil {
		return err
	}

	s.strategy = strategy
	return nil
}

func (s *DefaultNotificationService) Notify(message string, recipient string) error {
	return s.strategy.Send(message, recipient)
}

func (s *DefaultNotificationService) GetCurrentChannel() string {
	return s.strategy.GetChannelName()
}
