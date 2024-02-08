package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-service.greenocean-0f99656d.westeurope.azurecontainerapps.io/",
	}
}
