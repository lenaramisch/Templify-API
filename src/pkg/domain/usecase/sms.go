package usecase

func (u *Usecase) SendSMS(toNumber string, messageBody string) error {
	return u.smsSender.SendSMS(toNumber, messageBody)
}
