package user_service

import (
	"cloudStoregeDemo/models"
	"errors"
)

type User struct {
	ID                      int
	Username                string
	Email                   string
	ClientRandomValue       string
	EncryptedMasterKey      string
	HashedAuthenticationKey string
}

func (u *User) Login() (bool, error) {

	success, err := models.Validate(map[string]interface{}{
		"Username":                u.Username,
		"Email":                   u.Email,
		"HashedAuthenticationKey": u.HashedAuthenticationKey,
	})
	if err != nil {
		return false, err
	}
	if success {
		err := models.Login(map[string]interface{}{
			"HashedAuthenticationKey": u.HashedAuthenticationKey,
		})
		if err != nil {
			return true, err
		}
	}
	return success, nil
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
