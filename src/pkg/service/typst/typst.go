package typst

type TypstConfig struct {
}

type TypstService struct {
	config TypstConfig
}

func NewTypstService(config TypstConfig) *TypstService {
	return &TypstService{
		config: config,
	}
}

// TODO Implement FillTemplPlaceholders func
func (ts *TypstService) FillTemplatePlaceholders(typstString string, placeholders map[string]string) (string, error) {
	return "", nil
}
