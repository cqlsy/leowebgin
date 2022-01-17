package leowebgin

import (
	"fmt"
	"github.com/cqlsy/leoutil"
	"github.com/gorilla/websocket"
	"net/http"
)

//  这是一个测试的 socket初始化
func initTestSocket() {
	manager := NewManager(
		nil,
		func(client *SocketClient, msg []byte) {
			//leolog.LogInfoDefault(string(msg))

			// manager.Clients 这里保存了所有的socket的链接.需要的时候在这里寻找
			// 发送消息
			client.SendMessage([]byte(fmt.Sprintf("service callback：%s", msg)))
		},
	)
	// socket
	var s *Engine
	s = &Engine{
		engine:  nil,
		log:     nil,
		runMode: "",
	}
	s.AddSocketClient("/socket", manager, nil)
}

func (engine *Engine) AddSocketClient(path string, manager *SocketManager, upgrader *websocket.Upgrader) {
	engine.Get(path, manager.defaultWsHandler(upgrader))
	go manager.start()
}

// 生成当前的管理实例
func NewManager(geId func(c *Context) string,
	onGetMessage func(client *SocketClient, msg []byte)) *SocketManager {
	manage := &SocketManager{
		register:   make(chan *SocketClient),
		unregister: make(chan *SocketClient),
		Clients:    make(map[string]*SocketClient),
	}
	if geId == nil {
		geId = func(c *Context) string {
			// 当用户没有自定义当前的socket链接ID,这里生成默认的.
			return leoutil.RandString(32, "socket")
		}
	}
	manage.generateID = geId
	manage.onGetMessage = onGetMessage
	return manage
}

// 注册log的方式
func (manager *SocketManager) InitLog(f func(str string)) {
	manager.log = f
}

// 发送消息的方法
func (c *SocketClient) SendMessage(message []byte) {
	if c.isSendClose {
		return
	}
	c.send <- message
}

// 把http链接升级WebSocket的链接
func (manager *SocketManager) defaultWsHandler(Upgrader *websocket.Upgrader) func(c *Context) {
	return func(context *Context) {
		c := context.context
		if Upgrader == nil {
			Upgrader = upGraderDefault()
		}
		conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			http.NotFound(c.Writer, c.Request)
			return
		}
		// create client for ws
		client := &SocketClient{
			ID:          manager.generateID(c),
			socket:      conn,
			send:        make(chan []byte),
			isSendClose: false,
		}
		manager.register <- client
		// 启动当前通道的读写
		go client.read(manager)
		go client.write(manager)
	}
}

// 设置是否允许跨域
func upGraderDefault() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// 服务注册，注销的通道
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

// 读取信息的通道
func (c *SocketClient) read(manager *SocketManager) {
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

// 写入数据的通道
func (c *SocketClient) write(manager *SocketManager) {
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
