package leowebgin

import (
	"fmt"
	"github.com/cqlsy/leoutil"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func (w *WebGin) AddSocketClient(path string, manager *SocketManager, upgrader *websocket.Upgrader) {
	w.Gin.GET(path, manager.WsHandler(upgrader))
	go manager.start()
}

// Client is a websocket client
type Client struct {
	ID          string
	socket      *websocket.Conn
	send        chan []byte
	isSendClose bool
	ExtData     string // used for indentify data
}

// SocketManager is a websocket manager
type SocketManager struct {
	Clients      map[string]*Client
	register     chan *Client
	unregister   chan *Client
	generateID   func(c *gin.Context) string
	onGetMessage func(client *Client, msg []byte)
	log          func(errStr string)
}

func NewManager(geId func(c *gin.Context) string,
	onGetMessage func(client *Client, msg []byte)) *SocketManager {
	manage := &SocketManager{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
	if geId == nil {
		geId = func(c *gin.Context) string {
			return leoutil.RandString(32, "socket")
		}
	}
	manage.generateID = geId
	manage.onGetMessage = onGetMessage
	return manage
}

func (manager *SocketManager) InitLog(f func(str string)) {
	manager.log = f
}

func (c *Client) SendMessage(message []byte) {
	if c.isSendClose {
		return
	}
	c.send <- message
}

func (manager *SocketManager) WsHandler(Upgrader *websocket.Upgrader) func(c *gin.Context) {
	return func(c *gin.Context) {
		if Upgrader == nil {
			Upgrader = upGraderDefault()
		}
		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			http.NotFound(c.Writer, c.Request)
			return
		}
		// create client for ws
		client := &Client{
			ID:          manager.generateID(c),
			socket:      conn,
			send:        make(chan []byte),
			isSendClose: false,
		}
		manager.register <- client
		go client.read(manager)
		go client.write(manager)
	}
}

func (c *Client) read(manager *SocketManager) {
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			// if read message err ,we disconnect
			manager.unregister <- c
			_ = c.socket.Close()
			break
		}
		// callback on get message
		if manager.onGetMessage != nil {
			manager.onGetMessage(c, message)
		}
	}
}

// send message to user
func (c *Client) write(manager *SocketManager) {
	defer func() {
		_ = c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil && manager.log != nil {
					manager.log(err.Error())
				}
				return
			}
			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil && manager.log != nil {
				manager.log(err.Error())
			}
		}
	}
}

func upGraderDefault() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// restart on error
func (manager *SocketManager) start() {
	defer func() {
		err := recover()
		if err != nil {
			manager.log(fmt.Sprintf("%v", err))
		}
		go manager.start()
	}()
	for {
		select {
		case conn := <-manager.register:
			// connect success
			manager.Clients[conn.ID] = conn
			//conn.SendMessage([]byte("client "))
		case conn := <-manager.unregister:
			// dis connect
			if _, ok := manager.Clients[conn.ID]; ok {
				conn.isSendClose = true
				close(conn.send)
				delete(manager.Clients, conn.ID)
			}
		}
	}
}
