package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type RelayMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
type EnrichedRelayMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Clients struct {
	mu      sync.Mutex
	clients map[string]map[string]*websocket.Conn
}

type Client struct {
	connection *websocket.Conn
	groupId    string
	connId     string
}

var clients Clients = Clients{clients: make(map[string]map[string]*websocket.Conn)}

func main() {
	// Initialize database connection
	_, err := InitDB()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	defer CloseDB()

	fmt.Println("Starting WebSocket relay server on port 8080...")
	http.HandleFunc("/", handleConnection)
	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Check if auth key exists and get associated group id
func checkAuth(r *http.Request) string {
	authKey := r.URL.Query().Get("auth")
	return GetGroupByKey(authKey)
}

// Store client connection
func storeConn(conn *websocket.Conn, groupId string) string {
	clients.mu.Lock()
	defer clients.mu.Unlock()
	_, ok := clients.clients[groupId]
	if !ok {
		clients.clients[groupId] = make(map[string]*websocket.Conn)
	}
	connId := uuid.New().String()
	clients.clients[groupId][connId] = conn
	return connId
}

// Remove client connection
func removeConn(groupId string, connId string) {
	clients.mu.Lock()
	defer clients.mu.Unlock()

	if _, ok := clients.clients[groupId]; !ok {
		return
	}
	if _, ok := clients.clients[groupId][connId]; !ok {
		return
	}
	delete(clients.clients[groupId], connId)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	// Retrieve data and create store client
	groupId := checkAuth(r)
	if groupId == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("Forbidden"))
		return
	}
	connId := storeConn(conn, groupId)
	client := Client{
		connection: conn,
		groupId:    groupId,
		connId:     connId,
	}

	// Set close handler
	conn.SetCloseHandler(func(code int, text string) error {
		removeConn(groupId, connId)
		return nil
	})

	for {
		relayMsg := RelayMessage{}
		err := conn.ReadJSON(&relayMsg)
		if err != nil {
			fmt.Println("Error reading from websocket:", err)
			break
		}
		handleData(&client, &relayMsg)
	}
}

func handleData(client *Client, data *RelayMessage) {
	fmt.Printf("Recieved data from client: %s\n", data.Event)
	enrichedData := EnrichedRelayMessage{
		Event: data.Event,
		Data:  data.Data,
	}
	broadcastData(client, &enrichedData)
}

func broadcastData(client *Client, data *EnrichedRelayMessage) {
	if _, ok := clients.clients[client.groupId]; !ok {
		return
	}
	for k, v := range clients.clients[client.groupId] {
		if k != client.connId {
			v.WriteJSON(data)
		}
	}
}

func handleMessage(conn *websocket.Conn, message []byte) {
	fmt.Printf("Received message: %s\n", message)
	conn.WriteMessage(websocket.TextMessage, []byte("hello"))
}
