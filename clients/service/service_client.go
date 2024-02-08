package service

type Service struct {
	baseURL string
}

func New() *Service {
	return &Service{
		baseURL: "https://temp-service.happysmoke-1e47799e.eastus.azurecontainerapps.io/",
	}
}
