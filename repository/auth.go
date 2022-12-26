package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
)

type (
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	IAuthRepository interface {
		CreateUser(user model.CreateUserRequest, is_seller bool) error
		GetUserByEmail(email string) (*model.AuthUserDetails, error)
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	AuthRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// CreateUser repository layer for executing command creating a user
func (ar *AuthRepository) CreateUser(user model.CreateUserRequest, is_seller bool) error {
	tx, err := ar.DB.Begin(ar.Context)
	if err != nil {
		ar.Logger.Error(fmt.Sprintf("Error 13: %s", err))
		return err
	}
	q := `INSERT INTO "users" (full_name,email,password,role_id) VALUES ($1,$2,$3,$4) RETURNING id`

	var userID uuid.UUID

	err = tx.QueryRow(ar.Context, q, user.FullName, user.Email, user.Password, user.RoleID).Scan(&userID)
	if err != nil {
		err = tx.Rollback(ar.Context)
		if err != nil {
			ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.QueryRow.Scan Rollback ERROR %v MSG %s", err, err.Error()))
			return err
		}

		ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.QueryRow Scan ERROR %v MSG %s", err, err.Error()))
		return err
	}

	q2 := `INSERT INTO "users_profile" (telepon,location_id,user_id) VALUES ($1,$2,$3)`
	_, err = tx.Exec(ar.Context, q2, user.Telepon, user.LocationID, userID)
	if err != nil {
		err = tx.Rollback(ar.Context)
		if err != nil {
			ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.Exec User Profile Rollback ERROR %v MSG %s", err, err.Error()))
			return err
		}

		ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.Exec User Profile ERROR %v MSG %s", err, err.Error()))
		return err
	}

	if is_seller {
		qSeller := `INSERT INTO "lapak" (name,location_id,user_id) VALUES ($1,$2,$3)`
		var lapakName = "Lapak " + user.FullName
		_, err = tx.Exec(ar.Context, qSeller, lapakName, user.LocationID, userID)
		if err != nil {
			err = tx.Rollback(ar.Context)
			if err != nil {
				ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.Exec Lapak Penjual Rollback ERROR %v MSG %s", err, err.Error()))
				return err
			}

			ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser.Exec Lapak Penjual ERROR %v MSG %s", err, err.Error()))
			return err
		}
	}

	err = tx.Commit(ar.Context)
	if err != nil {
		ar.Logger.Error(fmt.Errorf("AuthRepository.CreateUser Commit ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}

// GetUserByEmail repository layer for querying command getting any user by email
func (ar *AuthRepository) GetUserByEmail(email string) (*model.AuthUserDetails, error) {
	var authUserDetail model.AuthUserDetails

	q := `SELECT 
		u.id AS user_id,
		u.full_name,
		u.email,
		u.password,
		r.id,
		r.name FROM "users" u 
		LEFT JOIN roles r ON r.id = u.role_id 
		WHERE u.email = $1
	`

	row := ar.DB.QueryRow(ar.Context, q, email)
	err := row.Scan(&authUserDetail.UserID, &authUserDetail.FullName, &authUserDetail.Email, &authUserDetail.Password, &authUserDetail.RoleID, &authUserDetail.RoleName)

	if err != nil {
		if err == pgx.ErrNoRows {
			ar.Logger.Info(fmt.Errorf("AuthRepository.GetUserByEmail INFO : %v MSG : %s", err, err.Error()))
		} else {
			ar.Logger.Error(fmt.Errorf("AuthRepository.GetUserByEmail ERROR : %v MSG : %s", err, err.Error()))
		}

		return nil, err
	}

	return &authUserDetail, nil
}
