package file_service

import "cloudStoregeDemo/models"

type File struct {
	models.Model
	ID           int
	ParentDictId int
	FileName     string
	EncryptedKey string
	FileContent  string
	FileType     string
	FileSize     int

	PageSize int
	PageNum  int
}

func (f *File) Add() (int, error) {
	return models.AddFile(map[string]interface{}{
		"ParentDictId": f.ParentDictId,
		"FileName":     f.FileName,
		"EncryptedKey": f.EncryptedKey,
		"FileContent":  f.FileContent,
		"FileType":     f.FileType,
		"FileSize":     f.FileSize,
	})
}

func (f *File) Get() (*models.FileSystem, string, error) {
	return models.GetFile(f.ID)
}

func (f *File) GetAll() ([]*models.FileSystem, error) {
	return models.GetFiles(f.ID, f.PageSize, f.PageNum)
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
