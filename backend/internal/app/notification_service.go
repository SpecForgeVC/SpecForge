package app

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type NotificationService interface {
	Register(client *websocket.Conn, userID uuid.UUID)
	Unregister(client *websocket.Conn)
	Broadcast(eventType string, payload interface{})
	NotifyUser(userID uuid.UUID, eventType string, payload interface{})
}

type notificationService struct {
	clients map[*websocket.Conn]uuid.UUID
	lock    sync.RWMutex
}

func NewNotificationService() NotificationService {
	return &notificationService{
		clients: make(map[*websocket.Conn]uuid.UUID),
	}
}

func (s *notificationService) Register(client *websocket.Conn, userID uuid.UUID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.clients[client] = userID
}

func (s *notificationService) Unregister(client *websocket.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.clients[client]; ok {
		delete(s.clients, client)
		client.Close()
	}
}

func (s *notificationService) Broadcast(eventType string, payload interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	msg := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	for client := range s.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			// If error, we might want to unregister, but let's avoid modifying map during iteration
			// In a real prod app, use a channel or separate go routine
			go client.Close()
		}
	}
}

func (s *notificationService) NotifyUser(userID uuid.UUID, eventType string, payload interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	msg := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	for client, uid := range s.clients {
		if uid == userID {
			client.WriteJSON(msg)
		}
	}
}
