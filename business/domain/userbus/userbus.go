package userbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/userbus/valueobject"
)

type User struct {
	ID          valueobject.ID
	FirstName   string
	LastName    string
	OtherNames  string
	Email       valueobject.Email
	Roles       []valueobject.Role
	Permissions []valueobject.Permission
	TenantID    string
	Tenant      valueobject.Tenant
	MFAEnrolled bool
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NewUser struct {
	FirstName  string
	LastName   string
	OtherNames string
	TenantID   string
	Email      valueobject.Email
}

type SearchResult struct {
	Users []User
	Total int64
}

type UpdateUser struct {
	FirstName  *string
	LastName   *string
	OtherNames *string
	Email      *valueobject.Email
}
