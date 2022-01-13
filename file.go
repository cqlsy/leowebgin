package leowebgin

import (
	"errors"
	"github.com/cqlsy/leofile"
	fileUtil "github.com/cqlsy/leofile/action"
	"github.com/gin-gonic/gin"
)

// Download the file, when the file is not directly accessible, you need to use this interface to access
func DownLoadFile(c *gin.Context, filePath string) error {
	if len(filePath) == 0 || !leofile.FileExists(filePath) {
		//c.JSON(200, ResponseMsg{-1, "error, no such file: " + filePath, ""})
		return errors.New("error, no such file: " + filePath)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileUtil.GetFilName(filePath))
	c.Header("Content-Transfer-Encoding", "binary")
	// down load
	//c.Header("Content-Disposition", "attachment; filename="+d["name"].(string))
	// show
	//c.Header("Content-Disposition", "inline;filename="+d["name"].(string))
	c.File(filePath)
	return nil
}

func UploadFile(c *gin.Context, key string) (string, error) {
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
func UploadFiles(c *gin.Context) ([]string, error) {
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
