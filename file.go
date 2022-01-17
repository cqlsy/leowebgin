package leowebgin

import (
	"errors"
	"github.com/cqlsy/leofile"
	fileUtil "github.com/cqlsy/leofile/action"
)

// Download the file, when the file is not directly accessible, you need to use this interface to access
// 下载文件时,调用这个方法
func DownLoadFile(context *Context, filePath string) error {
	c := context.context
	if len(filePath) == 0 || !leofile.FileExists(filePath) {
		//c.JSON(200, ResponseMsg{-1, "error, no such file: " + filePath, ""})
		return errors.New("error, no such file: " + filePath)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	// downLoad
	c.Header("Content-Disposition", "attachment; filename="+fileUtil.GetFilName(filePath))
	// show to explorer
	//c.Header("Content-Disposition", "inline;filename="+d["name"].(string))
	c.File(filePath)
	return nil
}

// 上传文件
func UploadFile(context *Context, key string) (string, error) {
	c := context.context
	fileHeader, err := c.FormFile(key)
	//fileHeader.Header.
	if err != nil {
		return "", err
	}
	fileReader, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer fileReader.Close()
	filPath, err := fileUtil.SaveFile(fileReader, fileHeader.Filename)
	return filPath, nil
}

// On many files
// 上传多个文件
func UploadFiles(context *Context) ([]string, error) {
	c := context.context
	var result []string
	form, err := c.MultipartForm()
	if err != nil {
		return result, errors.New("Get File err: " + err.Error())
	}
	for key, files := range form.File {
		_ = key
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				return nil, err
			}
			//file.Header.Get()
			filPath, err := fileUtil.SaveFile(f, file.Filename)
			if err != nil {
				return nil, err
			}
			result = append(result, filPath)
			// close on every time
			_ = f.Close()
		}
	}
	return result, nil
}
