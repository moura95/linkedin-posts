package main

import (
	"fmt"
	"log"
)

func main() {
	factory := &notificationFactory{}

	service, err := NewNotificationService(factory, "email")
	if err != nil {
		log.Fatalf("Error creating notification service: %v", err)
	}

	recipient := "usuario@exemplo.com"
	message := "Message!"

	fmt.Printf("Channel: %s\n", service.GetCurrentChannel())
	if err := service.Notify(message, recipient); err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	// Set new Strategy to SMS
	if err := service.SetStrategy("sms"); err != nil {
		log.Fatalf("Error to set SMS strategy: %v", err)
	}

	// Send SMS
	smsRecipient := "+5511987654321"
	fmt.Printf("\nChannel: %s\n", service.GetCurrentChannel())
	if err := service.Notify(message, smsRecipient); err != nil {
		log.Printf("Error setting SMS: %v", err)
	}

	// Set new Strategy to Push
	if err := service.SetStrategy("push"); err != nil {
		log.Fatalf("Error setting Push strategy: %v", err)
	}

	// Send Push
	deviceID := "device-xyz-123"
	fmt.Printf("\nChannel: %s\n", service.GetCurrentChannel())
	if err := service.Notify(message, deviceID); err != nil {
		log.Printf("Error sending Push: %v", err)
	}

	// Set new Strategy to Slack
	if err := service.SetStrategy("slack"); err != nil {
		log.Fatalf("Error setting Slack strategy: %v", err)
	}

	// Send Slack
	channel := "#geral"
	fmt.Printf("\nChannel: %s\n", service.GetCurrentChannel())
	if err := service.Notify(message, channel); err != nil {
		log.Printf("Error sending Slack: %v", err)
	}

	fmt.Println("\nNotificações enviadas com sucesso!")
}
