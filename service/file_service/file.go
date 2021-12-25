package file_service

import (
	"cloudStoregeDemo/models"
)

type File struct {
	models.Model
	ID           int
	UserId       int
	ParentDictId int
	FileName     string
	EncryptedKey string
	FileType     string
	FileSize     int

	PageSize int
	PageNum  int
}

func (f *File) Add() (int, error) {
	return models.AddFile(map[string]interface{}{
		"ParentDictId": f.ParentDictId,
		"FileName":     f.FileName,
		"UserId":       f.UserId,
		"EncryptedKey": f.EncryptedKey,
		"FileType":     f.FileType,
		"FileSize":     f.FileSize,
	})

	//if f.FileType == "-" {
	//	fileUrl := constant.FILE_SAVE_ROOT + strconv.Itoa(id) + "-" + f.FileContent.Filename
	//	f, err := os.Create(fileUrl)
	//	if err != nil {
	//		return 0, errors.New("文件创建失败")
	//	}
	//	_, err = f.WriteString(m["FileContent"].(string))
	//	if err != nil {
	//		return 0, errors.New("文件写入失败")
	//	}
	//}
}

func (f *File) Get() (*models.FileSystem, error) {
	return models.GetFile(f.ID, f.UserId)
}

func (f *File) GetAll() ([]*models.FileSystem, error) {
	return models.GetFiles(f.ID, f.UserId, f.PageSize, f.PageNum)
}

func (f *File) GetParentId() (int, error) {
	return models.GetParentFileId(f.ID)
}

func (f *File) Count() (int, error) {
	return models.GetFilesNum(f.ID)
}

func (f *File) Update() error {
	return models.UpdateFile(map[string]interface{}{
		"ID":           f.ID,
		"ParentDictId": f.ParentDictId,
		"FileName":     f.FileName,
		"FileSize":     f.FileSize,
	})
}

func (f *File) Delete() error {
	return models.DeleteFile(f.ID)
}
