package v1

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/app"
	"cloudStoregeDemo/pkg/constant"
	"cloudStoregeDemo/pkg/e"
	"cloudStoregeDemo/service/file_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
)

type File struct {
	models.Model
	ID           int    `json:"id"`
	ParentDictId int    `json:"parent_dict_id"`
	FileName     string `json:"file_name"`
	EncryptedKey string `json:"encrypted_key"`

	FileType string `json:"file_type"`
	FileSize int    `json:"file_size"`
	PageNum  int    `json:"page_num"`
	PageSize int    `json:"page_size"`
}

type FileForm struct {
	UploadFile   *multipart.FileHeader `form:"upload_file"`
	ParentDictId int                   `form:"parent_dict_id"`
	FileName     string                `form:"file_name"`
	EncryptedKey string                `form:"encrypted_key"`
	FileType     string                `form:"file_type"`
	FileSize     int                   `form:"file_size"`
}

// @Summary addfile
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /file/add [post]
func AddFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	var file = &FileForm{}
	var err error
	httpcode, errorcode := app.BindAndValid(c, file)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
		return
	}

	fmt.Println(file)
	userId, _ := appG.C.Get("userId")
	u, _ := strconv.Atoi(userId.(string))
	fileService := file_service.File{
		ParentDictId: file.ParentDictId,
		UserId:       u,
		EncryptedKey: file.EncryptedKey,
		FileType:     file.FileType,
		FileSize:     file.FileSize,
	}
	if file.FileType == "-" {
		fileService.FileName = file.UploadFile.Filename
		err = appG.C.SaveUploadedFile(file.UploadFile, constant.FILE_SAVE_ROOT+file.UploadFile.Filename)
		if err != nil {
			appG.Respond(http.StatusOK, e.ERROR, nil)
			return
		}
	} else {
		fileService.FileName = file.FileName
	}

	id, err := fileService.Add()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	if file.FileType == "-" {
		err = appG.C.SaveUploadedFile(file.UploadFile, constant.FILE_SAVE_ROOT+file.UploadFile.Filename)
		if err != nil {
			appG.Respond(http.StatusOK, e.ERROR, nil)
			return
		}
	}
	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"id": id,
	})

}

// @Summary getfiles
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /file/getAll [post]
func GetFiles(c *gin.Context) {
	var appG = app.Gin{C: c}
	var file = &File{}
	var err error
	httpcode, errorcode := app.BindAndValid(c, file)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
		return
	}
	userId, _ := appG.C.Get("userId")
	u, _ := strconv.Atoi(userId.(string))
	fileService := file_service.File{
		ID:       file.ID,
		UserId:   u,
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
	parentId, err := fileService.GetParentId()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	m := make(map[string]interface{})
	m["total"] = count
	m["files"] = files
	m["parent_dict_id"] = parentId

	appG.Respond(http.StatusOK, e.SUCCESS, m)
}

// @Summary getfile
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /file/get [get]
func GetFile(c *gin.Context) {
	var appG = app.Gin{C: c}
	id, err := strconv.Atoi(appG.C.Param("id"))
	if err != nil {
		appG.Respond(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	userId, _ := appG.C.Get("userId")
	u, _ := strconv.Atoi(userId.(string))
	fileService := file_service.File{
		ID:     id,
		UserId: u,
	}
	f, err := fileService.Get()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	fmt.Println(f.FileName)
	var HttpContentType = map[string]string{
		".avi":  "video/avi",
		".mp3":  "audio/mp3",
		".mp4":  "video/mp4",
		".wmv":  "video/x-ms-wmv",
		".asf":  "video/x-ms-asf",
		".rm":   "application/vnd.rn-realmedia",
		".rmvb": "application/vnd.rn-realmedia-vbr",
		".mov":  "video/quicktime",
		".m4v":  "video/mp4",
		".flv":  "video/x-flv",
		".jpg":  "image/jpeg",
		".png":  "image/png",
		".pdf":  "application/pdf",
		".docx": "application/msword",
		".doc":  "application/msword",
	}
	filePath := constant.FILE_SAVE_ROOT + f.FileName
	//获取文件名称带后缀
	fileNameWithSuffix := path.Base(filePath)
	//获取文件的后缀
	fileType := path.Ext(fileNameWithSuffix)
	//获取文件类型对应的http ContentType 类型
	fileContentType := HttpContentType[fileType]
	if fileContentType == "" {
		fileContentType = "application/octet-stream"
	}
	c.Header("Content-Type", fileContentType)
	c.File(filePath)

}

// @Summary updatefile
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /file/update [post]
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

// @Summary deletefile
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /file/delete [delete]
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
