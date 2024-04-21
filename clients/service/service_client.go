package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-service.icysky-10e92f2c.westeurope.azurecontainerapps.io/",
	}
}
