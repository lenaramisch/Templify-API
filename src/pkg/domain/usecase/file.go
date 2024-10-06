package usecase

import domain "templify/pkg/domain/model"

func (u *Usecase) GetFileUploadURL(fileName string) (*domain.FileUploadResponse, error) {
	//TODO check if file already exists
	return u.filemanagerService.GetFileUploadURL(fileName)
}

func (u *Usecase) GetFileDownloadURL(fileName string) (string, error) {
	return u.filemanagerService.GetFileDownloadURL(fileName)
}
