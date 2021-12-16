package models

import (
	"errors"
	"fmt"
	"time"
)

type FileSystem struct {
	Id           int `json:"id"`
	ParentDictId int
	FileName     string
	Ctime        time.Time
	mtime        time.Time
	atime        time.Time
	FileType     string
	FileSize     int
}

func (FileSystem) TableName() string {
	return "file_system"
}

func GetRootFileId(username string) {
	var err error
	var rootId int
	err = db.Select("root_dict_id").Where("username = ?", username).Find(&rootId).Error
	if err != nil {
		fmt.Printf("查询错误: %#v", err)
	}
	fmt.Println(rootId)
}

func AddFile(m map[string]interface{}) error {
	file := FileSystem{
		ParentDictId: m["ParentDictId"].(int),
		FileName:     m["FileName"].(string),
		Ctime:        time.Now(),
		FileType:     m["FileType"].(string),
		FileSize:     m["FileSize"].(int),
	}

	err := db.Save(&file).Error
	if err != nil {
		return errors.New("保存失败")
	}
	return nil
}

func UpdateFile(m map[string]interface{}) error {
	var err error
	var file = FileSystem{}
	err = db.Where("id = ?", m["ID"]).Take(&file).Error
	var flag int = 0
	if m["ParentDictId"] != "" && m["ParentDictId"] != file.ParentDictId {
		file.ParentDictId = m["ParentDictId"].(int)
		flag = 1
	}
	if m["FileName"] != "" && m["FileName"] != file.FileName {
		file.FileName = m["FileName"].(string)
		flag = 1
	}
	if m["FileSize"] != "" && m["FileSize"] != file.FileSize {
		file.FileSize = m["FileSize"].(int)
		flag = 1
	}
	file.mtime = time.Now()
	file.atime = time.Now()
	if flag == 1 {
		err = db.Save(file).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteFile(id int) error {
	return db.Where("id = ?", id).Delete(&FileSystem{}).Error
}

func GetFiles(id, pageSize, pageNum int) ([]*FileSystem, error) {
	var files []*FileSystem
	err := db.Where("parent_dict_id = ?", id).Offset(pageNum).Limit(pageSize).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetFile(id int) (*FileSystem, error) {
	var file FileSystem
	err := db.Where("id = ?", id).Find(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func GetFilesNum(id int) (int, error) {
	var count int
	err := db.Model(&FileSystem{}).Where("parent_dict_id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
