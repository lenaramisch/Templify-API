package mjmlservice

type MJMLConfig struct {
	//no config variables for now
}

type MJMLService struct {
	config MJMLConfig
}

func NewMJMLService(config MJMLConfig) *MJMLService {
	return &MJMLService{
		config: config,
	}
}

//TODO add MJMLService functions
