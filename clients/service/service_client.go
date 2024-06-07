package service

import "github.com/dro14/yordamchi/clients/other"

type Service struct {
	baseURL string
	apis    *other.APIs
}

func New() *Service {
	return &Service{
		baseURL: "https://yordamchi-service.greensmoke-1e04616b.westeurope.azurecontainerapps.io/",
		apis:    other.New(),
	}
}
