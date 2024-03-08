package web_socket

import (
	"awesomeProject/log_manager"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients     Clients
	fileManager *log_manager.FileManager
	mu          sync.Mutex
}

func NewManager() *Manager {
	manager := &Manager{
		clients:     make(Clients),
		fileManager: log_manager.GetFileManager(),
		mu:          sync.Mutex{},
	}
	go manager.SendUpdatedLog()
	return manager
}

func (m *Manager) addClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client]; ok {
		err := client.conn.Close()
		if err != nil {
			log.Println("Failed to close connection for remote add: ", client.conn.RemoteAddr())
			return
		}
		delete(m.clients, client)
	}
}

func (m *Manager) ServerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to initialize websocket for client having err: ", err)
	}
	client := NewClient(conn)
	m.addClient(client)
}

func (m *Manager) SendMessageToClient(event Event) {
	for client := range m.clients {
		if err := client.SendMessage(event); err != nil {
			m.removeClient(client)
		}
	}
}

func (m *Manager) readUpdatedLog() string {
	return m.fileManager.ReadLastLine()
}

func (m *Manager) SendUpdatedLog() {
	for {
		lastUpdatedLog := m.readUpdatedLog()
		if lastUpdatedLog != "" {
			payload, _ := json.Marshal(lastUpdatedLog)
			event := NewEvent(websocket.TextMessage, payload)
			m.SendMessageToClient(event)
		}
	}
}
