package multi_agent

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeQuestion   MessageType = "question"
	MessageTypeAnswer     MessageType = "answer"
	MessageTypeDelegation MessageType = "delegation"
	MessageTypeRefusal    MessageType = "refusal"
)

type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	Content   string      `json:"content"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Timestamp time.Time   `json:"timestamp"`
	ThreadID  string      `json:"thread_id"`
}

func NewMessage(msgType MessageType, content, from, to, threadID string) Message {
	return Message{
		ID:        uuid.New().String(),
		Type:      msgType,
		Content:   content,
		From:      from,
		To:        to,
		Timestamp: time.Now(),
		ThreadID:  threadID,
	}
}

type Conversation struct {
	ID         string    `json:"id"`
	Query      string    `json:"query"`
	Messages   []Message `json:"messages"`
	CreatedAt  time.Time `json:"created_at"`
	IsComplete bool      `json:"is_complete"`
}

func NewConversation(query string) *Conversation {
	return &Conversation{
		ID:         uuid.New().String(),
		Query:      query,
		Messages:   []Message{},
		CreatedAt:  time.Now(),
		IsComplete: false,
	}
}

func (c *Conversation) AddMessage(msg Message) {
	c.Messages = append(c.Messages, msg)
}
