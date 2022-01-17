package leowebgin

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
