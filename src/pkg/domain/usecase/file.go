package usecase

import domain "templify/pkg/domain/model"

func (u *Usecase) UploadFile(fileUploadRequest domain.FileUploadRequest) error {
	return u.fileManagerService.UploadFile(fileUploadRequest)
}

func (u *Usecase) DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error) {
	return u.fileManagerService.DownloadFile(fileDownloadRequest)
}
