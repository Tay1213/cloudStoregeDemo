package v1

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/app"
	"cloudStoregeDemo/pkg/e"
	"cloudStoregeDemo/service/file_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type File struct {
	models.Model
	ID           int    `json:"id"`
	ParentDictId int    `json:"parent_dict_id"`
	FileName     string `json:"file_name"`
	EncryptedKey string `json:"encrypted_key"`
	FileContent  string `json:"file_content"`

	FileType string `json:"file_type"`
	FileSize int    `json:"file_size"`
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
}

func AddFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	var file = &File{}
	//var err error
	httpcode, errorcode := app.BindAndValid(c, file)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
		return
	}

	fileService := file_service.File{
		ParentDictId: file.ParentDictId,
		FileName:     file.FileName,
		EncryptedKey: file.EncryptedKey,
		FileContent:  file.FileContent,
		FileType:     file.FileType,
		FileSize:     file.FileSize,
	}

	id, err := fileService.Add()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"id": id,
	})
}

func GetFiles(c *gin.Context) {
	var appG = app.Gin{C: c}
	var file = &File{}
	var err error
	httpcode, errorcode := app.BindAndValid(c, file)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
		return
	}

	fileService := file_service.File{
		ID:       file.ID,
		PageNum:  (file.PageNum - 1) * file.PageSize,
		PageSize: file.PageSize,
	}

	count, err := fileService.Count()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}

	files, err := fileService.GetAll()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	m := make(map[string]interface{})
	m["total"] = count
	m["files"] = files

	appG.Respond(http.StatusOK, e.SUCCESS, m)
}

func GetFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, err := strconv.Atoi(appG.C.Param("id"))
	if err != nil {
		appG.Respond(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	fileService := file_service.File{
		ID: id,
	}
	f, content, err := fileService.Get()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}

	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"file_msg":     f,
		"file_content": content,
	})
}

func UpdateFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	var file = &File{}
	httpcode, errorcode := app.BindAndValid(c, file)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
		return
	}

	fileService := file_service.File{
		ID:           file.ID,
		ParentDictId: file.ParentDictId,
		FileName:     file.FileName,
		FileType:     file.FileType,
		FileSize:     file.FileSize,
	}

	err := fileService.Update()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Respond(http.StatusOK, e.SUCCESS, nil)
}

func DeleteFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, err := strconv.Atoi(appG.C.Param("id"))
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	fileService := file_service.File{
		ID: id,
	}
	err = fileService.Delete()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Respond(http.StatusOK, e.SUCCESS, nil)
}
