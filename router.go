package leowebgin

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"net/http"
)

type Fun func(content *Context)

// 自定义 方法转化为 gin.HandlerFunc
func parseFun(fun ...Fun) []gin.HandlerFunc {
	var ff []gin.HandlerFunc
	for _, f := range fun {
		ff = append(ff, func(context *gin.Context) {
			f(&Context{context: context})
		})
	}
	return ff
}

func (engine *Engine) Group(path string, fun ...Fun) *RouterGroup {
	group := engine.engine.Group(path, parseFun(fun...)...)
	return &RouterGroup{group: group}
}

func (engine *RouterGroup) Get(relativePath string, fun ...Fun) {
	engine.group.GET(relativePath, parseFun(fun...)...)
}

func (engine *RouterGroup) Post(relativePath string, fun ...Fun) {
	engine.group.POST(relativePath, parseFun(fun...)...)
}

func (engine *Engine) Get(relativePath string, fun ...Fun) {
	engine.engine.GET(relativePath, parseFun(fun...)...)
}

func (engine *Engine) Post(relativePath string, fun ...Fun) {
	engine.engine.POST(relativePath, parseFun(fun...)...)
}

/*if path == "/" {
fmt.Println("you can't set GET Request when you set static path with '/'. " +
"Are you sure do this? ")
}*/
func (engine *Engine) StaticDir(path, filePath string) {
	engine.engine.StaticFS(path, http.Dir(filePath)) // 静态文件夹
	// engine.engine.StaticFile() 静态文件，
}

func (engine *Engine) StaticDirPackr(path, filePath string) {
	engine.engine.StaticFS(path, packr.New(path, filePath)) // 静态文件夹
	// engine.engine.StaticFile() 静态文件，
}
