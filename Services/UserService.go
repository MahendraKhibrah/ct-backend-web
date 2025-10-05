package Services

import (
	"ct-backend/Model"
	"ct-backend/Repository"
)

type (
	IUserService interface {
		GetAllVerified() (users []Model.User, err error)
	}
	UserService struct {
		UserRepository Repository.IUserRepository
	}
)

func UserServiceProvider(userRepository Repository.IUserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (h *UserService) GetAllVerified() (users []Model.User, err error) {
	return h.UserRepository.GetAllVerified()
}
