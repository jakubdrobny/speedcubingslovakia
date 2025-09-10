package interfaces

import (
	"context"

	"github.com/jakubdrobny/speedcubingslovakia/backend/models"
	"github.com/jakubdrobny/speedcubingslovakia/backend/shared/views"
)

type Repositories struct {
	userRepository UserRepository
}

type UserRepository interface {
	Update(ctx context.Context, db models.DbExecutor, user *models.User) error
	ViewManageUsers(ctx context.Context, db models.DbExecutor) ([]views.ManageUser, error)
}
