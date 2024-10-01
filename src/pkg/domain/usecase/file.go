package usecase

func (u *Usecase) GetFileUploadURL(fileName string) (string, error) {
	//TODO check if file already exists
	return u.fileManagerService.GetFileUploadURL(fileName)
}

func (u *Usecase) GetFileDownloadURL(fileName string) (string, error) {
	return u.fileManagerService.GetFileDownloadURL(fileName)
}
