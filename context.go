package leowebgin

import "io/ioutil"

/**
query 参数
*/
func (context *Context) GetQueryParams(name string) string {
	c := context.context
	params := c.Query(name)
	if params == "" {
		params = c.PostForm(name)
	}
	return params
}

func (content *Context) Abort() {
	content.context.Abort()
}

func (c *Context) IsAborted() bool {
	return c.context.IsAborted()
}

func (c *Context) Next() {
	c.context.Next()
}

func (c *Context) Header(key, value string) {
	c.context.Header(key, value)
}

// 获取 header 的数据
func (c *Context) GetHeader(key string) string {
	return c.context.GetHeader(key)
}
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.context.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *Context) Cookie(name string) (string, error) {
	return c.context.Cookie(name)
}

//router.GET("/user/:id", func(c *gin.Context) {
//    // a GET request to /user/john
//    id := c.Param("id") // id == "john"
//})
func (context *Context) GetQueryPath(name string) string {
	c := context.context
	params := c.Param(name)
	return params
}

func (context *Context) Query(key string) string {
	return context.context.Query(key)
}

/**
body 参数
*/
func (context *Context) GetBodyParams() (map[string]interface{}, error) {
	c := context.context
	data := make(map[string]interface{})
	err := c.ShouldBind(&data)
	return data, err
}

//params is a struct address
func (context *Context) GetBodyParamsStruct(params interface{}) error {
	c := context.context
	err := c.ShouldBind(params)
	return err
}

func (context *Context) ResponseJSON(obj interface{}) {
	context.context.JSON(200, obj)
}

func (context *Context) ResponseFile(filepath string) {
	context.context.File(filepath)
}

// 获取Body的数据
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.context.Request.Body)
}
