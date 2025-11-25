package service

import "personal-web-platform/internal/repository"

type Services struct {
	Profile ProfileService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Profile: NewProfileService(repos.Profile),
	}
}
