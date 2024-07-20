package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-service.greensmoke-1e04616b.westeurope.azurecontainerapps.io/",
	}
}
