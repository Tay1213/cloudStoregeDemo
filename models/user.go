package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	Id                      int
	Username                string
	Email                   string
	ClientRandomValue       string
	EncryptedMasterKey      string
	HashedAuthenticationKey string
	Regdate                 time.Time
	Logins                  int
	RootDictId              int
}

func (User) TableName() string {
	return "user"
}

func Validate(m map[string]interface{}) (bool, error) {
	var user = User{}
	if m["Username"] != "" {
		err := db.Where("username = ? and hashed_authentication_key = ?", m["Username"], m["HashedAuthenticationKey"]).
			Find(&user).Error
		if err != nil {
			return false, err
		}
	} else {
		err := db.Where("email = ? and hashed_authentication_key = ?", m["Email"], m["HashedAuthenticationKey"]).
			Find(&user).Error
		if err != nil {
			return false, err
		}
	}
	if user.Id != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func Login(m map[string]interface{}) (*User, error) {
	var user = &User{}
	var err error
	err = db.Where("hashed_authentication_key = ?", m["HashedAuthenticationKey"]).Find(&user).Error
	if err != nil {
		return nil, err
	}
	err = db.Model(&User{}).
		Where("hashed_authentication_key = ?", m["HashedAuthenticationKey"]).
		Update("logins", gorm.Expr("logins+1")).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetUserByName(username string) (*User, error) {
	var user = &User{}
	err := db.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func GetUserByEmail(email string) (*User, error) {
	var user = &User{}
	err := db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, err
}

func AddUser(m map[string]interface{}) error {
	var err error
	file := FileSystem{
		ParentDictId: 0,
		FileName:     "/",
		Ctime:        time.Now(),
		FileType:     "d",
		FileSize:     0,
	}

	err = db.Create(&file).Error
	if err != nil {
		return err
	}

	user := User{
		Email:                   m["Email"].(string),
		Username:                m["Username"].(string),
		ClientRandomValue:       m["ClientRandomValue"].(string),
		EncryptedMasterKey:      m["EncryptedMasterKey"].(string),
		HashedAuthenticationKey: m["HashedAuthenticationKey"].(string),
		Regdate:                 time.Now(),
		Logins:                  0,
		RootDictId:              file.Id,
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(m map[string]interface{}) error {
	var user = User{}
	var err error
	err = db.Where("id = ?", m["ID"]).Take(&user).Error
	if err != nil {
		return err
	}

	if m["Email"] != "" {
		err = db.Where("email = ?", m["Email"]).Find(&User{}).Error
		if err == nil {
			return errors.New("邮箱已存在")
		}
		user.Email = m["Email"].(string)
	}

	if m["Username"] != "" {
		err = db.Where("username = ?", m["Username"]).Find(&User{}).Error
		if err == nil {
			return errors.New("用户名已存在")
		}
		user.Username = m["Username"].(string)
	}
	err = db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil

}

func DeleteUser(id int) error {
	err := db.Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}

	return nil
}
