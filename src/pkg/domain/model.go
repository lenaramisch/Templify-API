package domain

type SmsRequest struct {
	ToNumber    string
	MessageBody string
}

type EmailRequest struct {
	ToEmail     string
	ToName      string
	Subject     string
	MessageBody string
}

type TemplateFillRequest struct {
	ShouldBeSent bool
	ToEmail      string
	ToName       string
	Subject      string
	Placeholders map[string]string
}

type EmailRequestAttm struct {
	ShouldBeSent bool
	ToEmail      string
	ToName       string
	Subject      string
	Placeholders map[string]string
	FileName     string
	FileType     string
	AttmContent  string
}

type Template struct {
	Name       string
	MJMLString string
}

type PDFTemplate struct {
	Name        string
	TypstString string
}
