package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-service.victoriousriver-fffd2d70.westeurope.azurecontainerapps.io/",
	}
}
