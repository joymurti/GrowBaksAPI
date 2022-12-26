package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
)

type (
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	IRoleRepository interface {
		GetRoleByName(name string) (*model.Role, error)
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	RoleRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// GetUserByEmail repository layer for querying command getting any user by email
func (rr *RoleRepository) GetRoleByName(name string) (*model.Role, error) {
	var role model.Role

	q := `SELECT 
		id,
		name
		FROM "roles"
		WHERE name = $1
	`

	row := rr.DB.QueryRow(rr.Context, q, name)
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			rr.Logger.Info(fmt.Errorf("RoleRepository.GetRoleByName INFO : %v MSG : %s", err, err.Error()))
		} else {
			rr.Logger.Error(fmt.Errorf("RoleRepository.GetRoleByName ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &role, nil
}
