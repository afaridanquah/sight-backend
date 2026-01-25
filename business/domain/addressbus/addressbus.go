package addressbus

import (
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/business/domain/addressbus/valueobject"
)

type Address struct {
	ID        valueobject.ID
	EntityID  string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   valueobject.Country
	Zip       string
	IsPrimary bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NewAddress struct {
	EntityID  string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   valueobject.Country
	IsPrimary bool
	Zip       string
}

type UpdateAddress struct {
	ID        valueobject.ID
	EntityID  string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   valueobject.Country
	Zip       string
	IsPrimary bool
}
