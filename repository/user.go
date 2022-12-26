package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/wiormiw/GrowBaks/config"
	"github.com/wiormiw/GrowBaks/model"
)

type (
	// IAuthRepository is an interface that has all the function to be implemented inside auth repository
	IUserRepository interface {
		GetAll(search string) ([]model.ViewUserResponse, error)
		GetByEmail(email string) (*model.ViewUserResponse, error)
		GetByID(id uuid.UUID) (*model.ViewUserResponse, error)
		GetProfileByID(id uuid.UUID) (*model.ViewUserProfileResponse, error)
		UpdateByID(id uuid.UUID, full_name string, email string, password string) error
		UpdateProfileByID(id uuid.UUID, email string, full_name string, tanggal_lahir string, gender string, telepon string) error
		DeleteByID(id uuid.UUID) error
	}

	// AuthRepository is an app auth struct that consists of all the dependencies needed for auth repository
	UserRepository struct {
		Context context.Context
		Config  *config.Configuration
		Logger  *logrus.Logger
		DB      *pgx.Conn
	}
)

// GetAll repository layer for querying command getting all user
func (ur *UserRepository) GetAll(search string) ([]model.ViewUserResponse, error) {
	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at
		FROM "users"
	`
	criteria := ""
	criteria = ""

	if len(search) > 0 {
		criteria += " full_name LIKE '%" + search + "%'"
	}

	cr := ""
	if len(criteria) > 0 {
		cr = " WHERE " + criteria
	}

	if cr != "" {
		q += cr
	}

	ur.Logger.Info(fmt.Sprintf("Query : %s", q))

	rows, err := ur.DB.Query(ur.Context, q)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.GetAll Query ERROR %v MSG %s", err, err.Error()))
		return nil, err
	}

	var listData []model.ViewUserResponse
	for rows.Next() {
		data := &model.ViewUserResponse{}
		err := rows.Scan(&data.ID, &data.FullName, &data.Email, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetAll rows.Next Scan ERROR %v MSG %s", err, err.Error()))
			return nil, err
		}

		listData = append(listData, *data)
	}

	return listData, nil
}

// GetByEmail repository layer for querying command get a user by email
func (ur *UserRepository) GetByEmail(email string) (*model.ViewUserResponse, error) {
	var user model.ViewUserResponse

	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at
		FROM "users"
		WHERE email = $1
	`

	row := ur.DB.QueryRow(ur.Context, q, email)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.Logger.Info(fmt.Errorf("UserRepository.GetByEmail Scan INFO %v MSG %s", err, err.Error()))
		} else {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetByEmail Scan ERROR %v MSG %s", err, err.Error()))
		}

		return nil, err
	}

	return &user, nil
}

// GetByID repository layer for querying command get a user by id
func (ur *UserRepository) GetByID(id uuid.UUID) (*model.ViewUserResponse, error) {
	var user model.ViewUserResponse

	q := `
		SELECT
			id,
			full_name,
			email,
			created_at,
			updated_at
		FROM "users"
		WHERE id = $1
	`

	row := ur.DB.QueryRow(ur.Context, q, id)
	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.Logger.Info(fmt.Errorf("UserRepository.GetByID Scan INFO %v MSG %s", err, err.Error()))
		} else {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetByID Scan ERROR %v MSG %s", err, err.Error()))
		}

		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetProfileByID(id uuid.UUID) (*model.ViewUserProfileResponse, error) {
	var userProfile model.ViewUserProfileResponse

	q := `SELECT u.id AS user_id,
		u.full_name,
		u.email,
		u.created_at,
		u.updated_at,
		up.telepon,
		up.gender,
		up.tanggal_lahir,
		loc.id as location_id,
		loc.daerah
		FROM "users_profile" up
		LEFT JOIN users u ON up.user_id = u.id
		LEFT JOIN "lokasi" loc ON loc.id = up.location_id
		WHERE user_id = $1
	`

	row := ur.DB.QueryRow(ur.Context, q, id)
	err := row.Scan(&userProfile.UserID,
		&userProfile.FullName,
		&userProfile.Email,
		&userProfile.CreatedAt,
		&userProfile.UpdatedAt,
		&userProfile.Telepon,
		&userProfile.Gender,
		&userProfile.TanggalLahir,
		&userProfile.LocationID,
		&userProfile.Daerah)
	if err != nil {
		if err == pgx.ErrNoRows {
			ur.Logger.Info(fmt.Errorf("UserRepository.GetProfileByID Scan INFO %v MSG %s", err, err.Error()))
		} else {
			ur.Logger.Error(fmt.Errorf("UserRepository.GetProfileByID Scan ERROR %v MSG %s", err, err.Error()))
		}

		return nil, err
	}

	return &userProfile, nil
}

// UpdateByID repository layer for executing command update a user by id
func (ur *UserRepository) UpdateByID(id uuid.UUID, full_name string, email string, password string) error {

	q := ` UPDATE users
		SET full_name = $1,
		email = $2,
		password = $3,
		updated_at = $4
	    WHERE id = $5
	`

	updatedAt := time.Now()

	_, err := ur.DB.Exec(ur.Context, q, full_name, email, password, updatedAt, id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.UpdateByID User ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateProfileByID(id uuid.UUID, email string, full_name string, tanggal_lahir string, gender string, telepon string) error {
	q1 := ` UPDATE users
		SET full_name = $1,
		email = $2,
		updated_at = $3
	    WHERE id = $4
	`

	_, err := ur.DB.Exec(ur.Context, q1, full_name, email, time.Now(), id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.UpdateProfileByID User Data ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	q2 := ` UPDATE users_profile
		SET telepon = $1,
		gender = $2,
		tanggal_lahir = $3,
		updated_at = $4
	    WHERE user_id = $5
	`
	_, err = ur.DB.Exec(ur.Context, q2, telepon, gender, tanggal_lahir, time.Now(), id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.UpdateProfileByID User Profile Data ERROR : %v MSG : %s", err, err.Error()))
		return err
	}

	return nil
}

// UpdateByID repository layer for executing command delete a user by id
func (ur *UserRepository) DeleteByID(id uuid.UUID) error {
	q := `DELETE FROM "users" WHERE id = $1`

	_, err := ur.DB.Exec(ur.Context, q, id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.DeleteByID Exec Users ERROR %v MSG %s", err, err.Error()))
		return err
	}

	q2 := `DELETE FROM "users_profile" WHERE user_id = $1`

	_, err = ur.DB.Exec(ur.Context, q2, id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.DeleteByID Exec Users Profile ERROR %v MSG %s", err, err.Error()))
		return err
	}

	q3 := `DELETE FROM "lapak" WHERE user_id = $1`

	_, err = ur.DB.Exec(ur.Context, q3, id)
	if err != nil {
		ur.Logger.Error(fmt.Errorf("UserRepository.DeleteByID Exec Lapak ERROR %v MSG %s", err, err.Error()))
		return err
	}

	return nil
}
