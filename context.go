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

//router.GET("/user/:id", func(c *gin.Context) {
//    // a GET request to /user/john
//    id := c.Param("id") // id == "john"
//})
func (context *Context) GetQueryPath(name string) string {
	c := context.context
	params := c.Param(name)
	return params
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

// 获取 header 的数据
func (c *Context) GetHeader(key string) string {
	return c.context.GetHeader(key)
}
