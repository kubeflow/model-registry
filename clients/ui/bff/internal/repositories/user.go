package repositories

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetUser(user string) (*models.User, error) {

	var formattedUser = user

	if formattedUser == "" {
		//if we are using token based auth, we still need to implement how to
		//safely get the user from the token
		formattedUser = "unknown"
	}
	var res = models.User{
		UserID: formattedUser,
	}

	return &res, nil
}
