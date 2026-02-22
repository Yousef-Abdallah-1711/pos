//go:build legacy
// +build legacy

package service

import (
	"fmt"
	"sync"
	"time"
)

// NotificationService handles sending notifications for alerts and events
type NotificationService struct {
	channels      map[string]NotificationChannel
	templates     map[string]*NotificationTemplate
	rules         map[string]*NotificationRule
	history       []*NotificationRecord
	mu            sync.RWMutex
	queue         chan *Notification
	maxQueueSize  int
	workers       int
}

// NotificationChannel defines a notification delivery channel
type NotificationChannel struct {
	ID            string
	Type          string // "email", "slack", "webhook", "sms"
	Name          string
	Configuration map[string]interface{}
	IsActive      bool
	CreatedAt     time.Time
}

// Notification represents a notification to send
type Notification struct {
	ID        string
	Type      string // "alert", "reminder", "notification"
	Title     string
	Message   string
	Severity  string // "info", "warning", "critical"
	Channels  []string
	Data      map[string]interface{}
	Timestamp time.Time
	Retries   int
}

// NotificationTemplate represents a message template
type NotificationTemplate struct {
	ID           string
	Name         string
	Subject      string
	Body         string
	HTMLBody     string
	Variables    []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NotificationRule defines when to send notifications
type NotificationRule struct {
	ID              string
	Name            string
	Description     string
	Condition       string // "alert_triggered", "sla_violated", "threshold_exceeded"
	ChannelID       string
	TemplateID      string
	Severity        string
	RateLimit       int // Max notifications per hour
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NotificationRecord tracks sent notifications
type NotificationRecord struct {
	ID        string
	NotificationID string
	ChannelID string
	Status    string // "pending", "sent", "failed"
	SentAt    *time.Time
	Error     string
	Attempt   int
	CreatedAt time.Time
}

// NewNotificationService creates a new notification service
func NewNotificationService(workers int) *NotificationService {
	return &NotificationService{
		channels:     make(map[string]NotificationChannel),
		templates:    make(map[string]*NotificationTemplate),
		rules:        make(map[string]*NotificationRule),
		history:      make([]*NotificationRecord, 0),
		queue:        make(chan *Notification, 1000),
		maxQueueSize: 10000,
		workers:      workers,
	}
}

// Start starts the notification service
func (ns *NotificationService) Start() {
	fmt.Println("DEBUG: Starting notification service")

	// Start worker goroutines
	for i := 0; i < ns.workers; i++ {
		go ns.worker(i)
	}
}

// SendNotification sends a notification
func (ns *NotificationService) SendNotification(notification *Notification) error {
	if notification.ID == "" {
		notification.ID = fmt.Sprintf("notif_%d", time.Now().UnixNano())
	}
	if notification.Timestamp.IsZero() {
		notification.Timestamp = time.Now()
	}

	// Add to queue
	select {
	case ns.queue <- notification:
		fmt.Printf("DEBUG: Notification queued: %s\n", notification.ID)
		return nil
	default:
		return fmt.Errorf("notification queue full")
	}
}

// AddChannel adds a notification channel
func (ns *NotificationService) AddChannel(channel *NotificationChannel) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if channel.ID == "" {
		channel.ID = fmt.Sprintf("channel_%d", time.Now().UnixNano())
	}
	channel.CreatedAt = time.Now()

	ns.channels[channel.ID] = *channel
	fmt.Printf("DEBUG: Notification channel added: %s (%s)\n", channel.Name, channel.Type)

	return nil
}

// GetChannel retrieves a notification channel
func (ns *NotificationService) GetChannel(channelID string) (*NotificationChannel, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	channel, exists := ns.channels[channelID]
	if !exists {
		return nil, fmt.Errorf("channel not found")
	}

	return &channel, nil
}

// ListChannels lists all notification channels
func (ns *NotificationService) ListChannels() []*NotificationChannel {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	channels := make([]*NotificationChannel, 0, len(ns.channels))
	for _, ch := range ns.channels {
		channels = append(channels, &ch)
	}

	return channels
}

// AddTemplate adds a notification template
func (ns *NotificationService) AddTemplate(template *NotificationTemplate) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if template.ID == "" {
		template.ID = fmt.Sprintf("template_%d", time.Now().UnixNano())
	}
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	ns.templates[template.ID] = template
	fmt.Printf("DEBUG: Notification template added: %s\n", template.Name)

	return nil
}

// GetTemplate retrieves a notification template
func (ns *NotificationService) GetTemplate(templateID string) (*NotificationTemplate, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	template, exists := ns.templates[templateID]
	if !exists {
		return nil, fmt.Errorf("template not found")
	}

	return template, nil
}

// AddRule adds a notification rule
func (ns *NotificationService) AddRule(rule *NotificationRule) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if rule.ID == "" {
		rule.ID = fmt.Sprintf("rule_%d", time.Now().UnixNano())
	}
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	ns.rules[rule.ID] = rule
	fmt.Printf("DEBUG: Notification rule added: %s\n", rule.Name)

	return nil
}

// GetRule retrieves a notification rule
func (ns *NotificationService) GetRule(ruleID string) (*NotificationRule, error) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	rule, exists := ns.rules[ruleID]
	if !exists {
		return nil, fmt.Errorf("rule not found")
	}

	return rule, nil
}

// GetNotificationHistory retrieves notification history
func (ns *NotificationService) GetNotificationHistory(limit int) []*NotificationRecord {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	if limit > len(ns.history) {
		limit = len(ns.history)
	}

	records := make([]*NotificationRecord, limit)
	copy(records, ns.history[len(ns.history)-limit:])

	return records
}

// Worker processes notifications
func (ns *NotificationService) worker(id int) {
	fmt.Printf("DEBUG: Notification worker %d started\n", id)

	for notification := range ns.queue {
		ns.processNotification(notification)
	}

	fmt.Printf("DEBUG: Notification worker %d stopped\n", id)
}

// ProcessNotification sends a notification through configured channels
func (ns *NotificationService) processNotification(notification *Notification) {
	fmt.Printf("DEBUG: Processing notification: %s\n", notification.ID)

	ns.mu.RLock()
	channels := notification.Channels
	ns.mu.RUnlock()

	for _, channelID := range channels {
		go ns.sendToChannel(notification, channelID)
	}
}

// SendToChannel sends notification to a specific channel
func (ns *NotificationService) sendToChannel(notification *Notification, channelID string) {
	channel, err := ns.GetChannel(channelID)
	if err != nil {
		fmt.Printf("DEBUG: Channel not found: %s\n", channelID)
		return
	}

	if !channel.IsActive {
		fmt.Printf("DEBUG: Channel inactive: %s\n", channelID)
		return
	}

	// Record attempt
	record := &NotificationRecord{
		ID:             fmt.Sprintf("record_%d", time.Now().UnixNano()),
		NotificationID: notification.ID,
		ChannelID:      channelID,
		Status:         "pending",
		Attempt:        notification.Retries + 1,
		CreatedAt:      time.Now(),
	}

	// Send based on channel type
	var sendErr error
	switch channel.Type {
	case "email":
		sendErr = ns.sendEmail(notification, channel)
	case "slack":
		sendErr = ns.sendSlack(notification, channel)
	case "webhook":
		sendErr = ns.sendWebhook(notification, channel)
	case "sms":
		sendErr = ns.sendSMS(notification, channel)
	default:
		sendErr = fmt.Errorf("unsupported channel type: %s", channel.Type)
	}

	if sendErr != nil {
		record.Status = "failed"
		record.Error = sendErr.Error()

		// Retry if needed
		if notification.Retries < 3 {
			notification.Retries++
			time.Sleep(5 * time.Second)
			ns.processNotification(notification)
		}
	} else {
		record.Status = "sent"
		now := time.Now()
		record.SentAt = &now
		fmt.Printf("DEBUG: Notification sent via %s: %s\n", channel.Type, notification.ID)
	}

	// Record in history
	ns.mu.Lock()
	ns.history = append(ns.history, record)
	ns.mu.Unlock()
}

// SendEmail sends notification via email (placeholder)
func (ns *NotificationService) sendEmail(notification *Notification, channel *NotificationChannel) error {
	fmt.Printf("DEBUG: Sending email notification: %s\n", notification.ID)
	// In production, would use SMTP or email service API
	return nil
}

// SendSlack sends notification via Slack (placeholder)
func (ns *NotificationService) sendSlack(notification *Notification, channel *NotificationChannel) error {
	fmt.Printf("DEBUG: Sending Slack notification: %s\n", notification.ID)
	// In production, would make HTTP request to Slack webhook
	return nil
}

// SendWebhook sends notification via webhook
func (ns *NotificationService) sendWebhook(notification *Notification, channel *NotificationChannel) error {
	fmt.Printf("DEBUG: Sending webhook notification: %s\n", notification.ID)
	// In production, would make HTTP request to configured webhook URL
	return nil
}

// SendSMS sends notification via SMS (placeholder)
func (ns *NotificationService) sendSMS(notification *Notification, channel *NotificationChannel) error {
	fmt.Printf("DEBUG: Sending SMS notification: %s\n", notification.ID)
	// In production, would use SMS service API (Twilio, etc.)
	return nil
}

// Shutdown gracefully shuts down the notification service
func (ns *NotificationService) Shutdown() error {
	fmt.Println("DEBUG: Shutting down notification service")
	close(ns.queue)
	return nil
}
