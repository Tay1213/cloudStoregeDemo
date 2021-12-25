package user_service

import (
	"cloudStoregeDemo/models"
	"cloudStoregeDemo/pkg/constant"
	"cloudStoregeDemo/pkg/gredis"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type User struct {
	ID                      int
	Username                string
	Email                   string
	ClientRandomValue       string
	EncryptedMasterKey      string
	HashedAuthenticationKey string
}

func (u *User) Login() (bool, string, error) {

	success, err := models.Validate(map[string]interface{}{
		"Username":                u.Username,
		"Email":                   u.Email,
		"HashedAuthenticationKey": u.HashedAuthenticationKey,
	})
	if err != nil || !success {
		return false, "", err
	}
	user, err := models.Login(map[string]interface{}{
		"HashedAuthenticationKey": u.HashedAuthenticationKey,
	})
	if err != nil {
		return false, "", err
	}

	flag, _ := gredis.Get(constant.LOGIN_PRIFIX + ":" + strconv.Itoa(user.Id))
	if flag != nil {
		times := flag[0] - 48
		if times >= 5 {
			return false, "", errors.New("登录次数过多")
		} else {
			err = gredis.Set(constant.LOGIN_PRIFIX+":"+strconv.Itoa(user.Id), times+1, 30)
			if err != nil {
				return false, "", errors.New("设置redis出错")
			}
		}
	} else {
		err = gredis.Set(constant.LOGIN_PRIFIX+":"+strconv.Itoa(user.Id), 1, 30)
		if err != nil {
			return false, "", errors.New("设置redis出错")
		}
	}
	token, err := getJwt(user)
	if err != nil {
		return false, "", errors.New("获取token失败")
	}
	return success, token, nil
}

func (u *User) GetUserByName() (*models.User, error) {
	return models.GetUserByName(u.Username)
}

func (u *User) GetUserByEmail() (*models.User, error) {
	return models.GetUserByEmail(u.Email)
}

func (u *User) Reg() error {
	var err error
	_, err = models.GetUserByEmail(u.Email)
	if err == nil {
		return errors.New("用户名或邮箱已被注册")
	}
	_, err = models.GetUserByName(u.Username)
	if err == nil {
		return errors.New("用户名或邮箱已被注册")
	}

	err = models.AddUser(map[string]interface{}{
		"Email":                   u.Email,
		"Username":                u.Username,
		"ClientRandomValue":       u.ClientRandomValue,
		"EncryptedMasterKey":      u.EncryptedMasterKey,
		"HashedAuthenticationKey": u.HashedAuthenticationKey,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Update() error {
	return models.UpdateUser(map[string]interface{}{
		"ID":       u.ID,
		"Email":    u.Email,
		"Username": u.Username,
	})
}

func (u *User) Delete() error {
	return models.DeleteUser(u.ID)
}

func getJwt(user *models.User) (string, error) {
	id := strconv.Itoa(user.Id)
	claims := jwt.StandardClaims{
		Audience:  user.Username,                                         // 受众
		ExpiresAt: time.Now().Unix() + (int64(constant.JWT_EXPIRE_TIME)), // 失效时间
		Id:        id,                                                    // 编号
		IssuedAt:  time.Now().Unix(),                                     // 签发时间
		Issuer:    "cloud_storage",                                       // 签发人
		NotBefore: time.Now().Unix(),                                     // 生效时间
		Subject:   "login",                                               // 主题
	}
	var jwtSecret = []byte(constant.JWT_SECRET)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	err = gredis.Set(constant.TOKEN_PRIFIX+":"+id+":"+token, "token", constant.EXPIRE_TIME)
	if err != nil {
		return "", err
	}
	return token, nil
}
