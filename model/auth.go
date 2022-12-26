package model

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	// KeyJWTValidAccess is context key identifier for valid jwt token
	KeyJWTValidAccess = "ValidJWTAccess"

	// IsAllowedEmailInput is regex validator to allowing only valid email
	IsAllowedEmailInput *regexp.Regexp

	// Regex Nomor Telepon
	IsAllowedTeleponInput *regexp.Regexp

	// Regex Tanggal
	IsAllowedDateInput *regexp.Regexp
)

type (
	// AuthDetails consist data authorized users
	AuthUserDetails struct {
		UserID   uuid.UUID `db:"user_id" json:"-"`
		FullName string    `db:"full_name" json:"full_name"`
		Email    string    `db:"email" json:"email"`
		RoleID   uuid.UUID `db:"id" json:"role_id"`
		RoleName string    `db:"name" json:"role_name"`
		Password string    `db:"password" json:"-"`
	}

	// CreateUserRequest consist data for creating a user
	CreateUserRequest struct {
		FullName   string    `json:"full_name"`
		Email      string    `json:"email"`
		Password   string    `json:"password"`
		Telepon    string    `json:"telepon"`
		LocationID string    `json:"location_id"`
		IsSeller   bool      `json:"is_seller"`
		RoleID     uuid.UUID `json:"-"`
	}

	// LoginRequest consist data for log-in a user
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// SuccessLoginResponse consist data of success login
	SuccessLoginResponse struct {
		AccessToken string    `json:"access_token"`
		ExpiredAt   time.Time `json:"expired_at"`
		Role        string    `json:"role"`
	}
)

const (
	// Role Mapping
	RoleSeller   string = "Penjual"
	RoleCustomer string = "Pembeli"
	RoleAdmin    string = "Admin"
)

func init() {
	IsAllowedEmailInput = regexp.MustCompile(`^[^\s@]+@([^\s@.,]+\.)+[^\s@.,]{2,}$`)
	IsAllowedTeleponInput = regexp.MustCompile(`^(\+62|62)?[\s-]?0?8[1-9]{1}\d{1}[\s-]?\d{4}[\s-]?\d{2,5}$`)
	IsAllowedDateInput = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
}
