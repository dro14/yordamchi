package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-python.greensmoke-1e04616b.westeurope.azurecontainerapps.io/",
	}
}
