package usecase

import (
	domain "templify/pkg/domain/model"
)

func (u *Usecase) GetFileUploadURL(fileName string) (*domain.FileUploadResponse, error) {
	// Check if file already exists with GetFileDownloadURL
	_, err := u.filemanagerService.GetFileDownloadURL(fileName)
	if err == nil {
		return nil, domain.ErrorFileAlreadyExists{FileName: fileName}
	}
	return u.filemanagerService.GetFileUploadURL(fileName)
}

func (u *Usecase) GetFileDownloadURL(fileName string) (string, error) {
	return u.filemanagerService.GetFileDownloadURL(fileName)
}
