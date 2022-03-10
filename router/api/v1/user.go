package v1

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/app"
	"cloudStoregeDemo/pkg/e"
	"cloudStoregeDemo/pkg/util"
	"cloudStoregeDemo/service/file_service"
	"cloudStoregeDemo/service/user_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type User struct {
	ID                      int    `json:"id"`
	Username                string `json:"username"`
	Email                   string `json:"email"`
	ClientRandomValue       string `json:"client_random_value"`
	HashedAuthenticationKey string `json:"hashed_authentication_key"`
	EncryptedMasterKey      string `json:"encrypted_master_key"`
}

// @Summary login
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/login [post]
func Login(c *gin.Context) {
	var appG = app.Gin{C: c}
	var user User

	httpCode, errCode := app.BindAndValid(c, &user)
	if errCode != e.SUCCESS {
		appG.Respond(httpCode, errCode, nil)
		return
	}
	fmt.Println(user)
	userService := user_service.User{
		Username:                user.Username,
		Email:                   user.Email,
		HashedAuthenticationKey: user.HashedAuthenticationKey,
	}
	success, token, err := userService.Login()
	if err != nil || !success {
		appG.RespondMsg(http.StatusOK, e.ERROR, err.Error(), nil)
		return
	}
	var u = &models.User{}
	u, err = userService.GetUserByName()
	if err != nil {
		u, _ = userService.GetUserByEmail()
	}

	fileService := file_service.File{
		ID:       u.RootDictId,
		UserId:   u.Id,
		PageNum:  0,
		PageSize: 10,
	}
	total, err := fileService.Count()
	files, err := fileService.GetAll()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, nil)
		return
	}

	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"id":                   u.Id,
		"login":                success,
		"token":                "Bearer " + token,
		"root_dict_id":         u.RootDictId,
		"encrypted_master_key": u.EncryptedMasterKey,
		"total":                total,
		"files":                files,
	})
}

// @Summary logout
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/logout [post]
func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Respond(http.StatusOK, e.SUCCESS, nil)
}

// @Summary GetUserByName
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/name [get]
func GetUserByName(c *gin.Context) {
	var appG = app.Gin{C: c}
	name := appG.C.Param("name")
	userService := user_service.User{
		Username: name,
	}

	user, err := userService.GetUserByName()
	if err != nil {
		//appG.Respond(http.StatusOK, e.USER_NOT_FOUND, nil)
		//return
		appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"client_random_value": util.RandStringBytesMaskImprSrc1(32),
		})
		return
	}

	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"client_random_value": user.ClientRandomValue,
	})
}

// @Summary GetUserByEmail
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/email [get]
func GetUserByEmail(c *gin.Context) {
	var appG = app.Gin{C: c}
	email := appG.C.Param("email")
	userService := user_service.User{
		Email: email,
	}

	user, err := userService.GetUserByEmail()
	if err != nil {
		//appG.Respond(http.StatusOK, e.EMAIL_NOT_FOUND, nil)
		//return
		appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
			"client_random_value": util.RandStringBytesMaskImprSrc1(32),
		})
	}

	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"client_random_value": user.ClientRandomValue,
	})
}

// @Summary register
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/reg [post]
func Reg(c *gin.Context) {
	var appG = app.Gin{C: c}
	var user = &User{}
	httpcode, errcode := app.BindAndValid(c, &user)
	if errcode != e.SUCCESS {
		appG.Respond(httpcode, errcode, nil)
		return
	}

	userService := user_service.User{
		Username:                user.Username,
		Email:                   user.Email,
		ClientRandomValue:       user.ClientRandomValue,
		EncryptedMasterKey:      user.EncryptedMasterKey,
		HashedAuthenticationKey: user.HashedAuthenticationKey,
	}

	err := userService.Reg()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}

	appG.Respond(http.StatusOK, e.SUCCESS, nil)

}

// @Summary UpdateUser
// @Produce  json
// @Param user body string false "username, email, hashedAuthenticationKey"
// @Success 200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/update [put]
func UpdateUser(c *gin.Context) {
	var appG = app.Gin{C: c}
	var user = &User{}
	httpcode, errorcode := app.BindAndValid(c, user)
	if errorcode != e.SUCCESS {
		appG.Respond(httpcode, errorcode, nil)
	}

	userService := user_service.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	err := userService.Update()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Respond(http.StatusOK, e.SUCCESS, nil)
}

// @Summary DeleteUser
// @Produce  json
// @Param  user body string false "username, email, hashedAuthenticationKey"
// @Success   200 {object} app.ResultData
// @Failure 500 {object} app.ResultData
// @Router /user/delete [delete]
func DeleteUser(c *gin.Context) {
	var appG = app.Gin{C: c}
	var err error
	id, err := strconv.Atoi(appG.C.Param("id"))
	if err != nil {
		appG.Respond(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	userService := user_service.User{
		ID: id,
	}

	err = userService.Delete()
	if err != nil {
		appG.Respond(http.StatusOK, e.ERROR, err.Error())
		return
	}
	appG.Respond(http.StatusOK, e.SUCCESS, nil)

}
