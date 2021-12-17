package v1

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/app"
	"cloudStoregeDemo/pkg/e"
	"cloudStoregeDemo/pkg/util"
	"cloudStoregeDemo/service/file_service"
	"cloudStoregeDemo/service/user_service"
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
	userService := user_service.User{
		Username:                user.Username,
		Email:                   user.Email,
		HashedAuthenticationKey: user.HashedAuthenticationKey,
	}
	success, err := userService.Login()
	if err != nil {
		appG.Respond(http.StatusOK, 501, nil)
		return
	}
	var u = &models.User{}
	u, err = userService.GetUserByName()
	if err != nil {
		u, _ = userService.GetUserByEmail()
	}

	fileService := file_service.File{
		ID:       u.RootDictId,
		PageNum:  0,
		PageSize: 10,
	}

	total, err := fileService.Count()

	files, err := fileService.GetAll()
	if err != nil {
		appG.Respond(http.StatusOK, 501, nil)
		return
	}

	appG.Respond(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"login":        success,
		"root_dict_id": u.RootDictId,
		"total":        total,
		"files":        files,
	})
}

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
