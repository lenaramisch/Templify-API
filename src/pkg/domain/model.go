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

type Template struct {
	Name       string
	MJMLString string
}
