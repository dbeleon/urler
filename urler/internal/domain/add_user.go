package domain

import (
	"context"

	"github.com/dbeleon/urler/urler/internal/domain/models"
)

func (m *Model) AddUser(ctx context.Context, name, email string) (int64, error) {
	u, err := m.repo.AddUser(models.User{Name: name, Email: email})
	if err != nil {
		return 0, err
	}

	return u.Id, nil
}
