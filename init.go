package leowebgin

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RunMode string

const (
	Pro RunMode = "pro"
	Dev RunMode = "dev"
)

type WebConf struct {
	Ip      string
	Port    int
	LogPath string
	RunMode string
	SaveFilePath string
	PicDoMain    string
}

// 日志打印
type Logger func(isPro bool, logStr string)

/**
gin的核心发动机
*/
type Engine struct {
	engine  *gin.Engine
	log     Logger
	runMode RunMode
}

// 组路由
type RouterGroup struct {
	group *gin.RouterGroup
}

/**
gin 的上下文Context
*/
type Context struct {
	context *gin.Context
}

/**
Socket 所需要的实例
*/
// SocketManager is a websocket manager
type SocketManager struct {
	Clients      map[string]*SocketClient               // 存储的所有的实例链接
	register     chan *SocketClient                     // 注册服务
	unregister   chan *SocketClient                     // 注销服务
	generateID   func(c *Context) string                // 生成唯一标识
	onGetMessage func(client *SocketClient, msg []byte) // 接收消息
	log          func(errStr string)                    // 日志输出
}

// 单个socket 链接实例
type SocketClient struct {
	ID          string          // 链接Id（唯一标识）
	socket      *websocket.Conn // websocket 链接实例
	send        chan []byte     // 当前链接的消息发送
	isSendClose bool            // 当前链接是否关闭
	ExtData     string          // used for indentify data

}
