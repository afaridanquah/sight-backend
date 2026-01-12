package valueobject

import (
	"errors"
	"strings"
	"time"
)

type Document struct {
	ID             string
	Parent         string
	DocumentType   IdentificationType
	Side           Side
	OriginalName   string
	FileName       string
	DocumentStatus DocumentStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func ParseDocument(
	id string,
	parent *string,
	docT IdentificationType,
	side Side,
	origFile string,
	fileName string,
	status DocumentStatus,
	createdAt time.Time,
	updatedAt time.Time) (Document, error) {
	if id == "" {
		return Document{}, errors.New("id is required")
	}

	doc := Document{
		ID:             id,
		Parent:         *parent,
		DocumentType:   docT,
		Side:           side,
		OriginalName:   origFile,
		FileName:       fileName,
		DocumentStatus: status,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	return doc, nil
}

// ====================================================

type Side struct {
	a string
}

var (
	ErrSideCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	FRONT_SIDE = Side{"FRONT"}
	BACK_SIDE  = Side{"BACK"}
)

var AvailableSides = []Side{FRONT_SIDE, BACK_SIDE}

func ParseSide(name string) (Side, error) {
	if name == "" {
		return Side{}, ErrSideCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "FRONT":
		return FRONT_SIDE, nil
	case "BACK":
		return BACK_SIDE, nil
	default:
		return Side{}, errors.New(" invalid name")
	}
}

func MustParseSide(name string) Side {
	status, err := ParseSide(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s Side) String() string {
	return s.a
}

// ============================================================
type DocumentStatus struct {
	a string
}

var (
	ErrStatusCannotBeEmpty = errors.New("status name cannot be empty")
)

var (
	VERIFIED = DocumentStatus{"VERIFIED"}
	REJECTED = DocumentStatus{"REJECTED"}
	PENDING  = DocumentStatus{"PENDING"}
	DRAFT    = DocumentStatus{"DRAFT"}
)

var Statuses = []DocumentStatus{VERIFIED, REJECTED}

func ParseStatus(name string) (DocumentStatus, error) {
	if name == "" {
		return DocumentStatus{}, ErrStatusCannotBeEmpty
	}
	name = strings.ToUpper(name)
	switch name {
	case "VERIFIED":
		return VERIFIED, nil
	case "REJECTED":
		return REJECTED, nil
	case "PENDING":
		return PENDING, nil
	case "DRAFT":
		return DRAFT, nil
	default:
		return DocumentStatus{}, errors.New("status name : invalid name")
	}
}

func MustParseStatus(name string) DocumentStatus {
	status, err := ParseStatus(name)
	if err != nil {
		panic(err)
	}
	return status
}

func (s DocumentStatus) String() string {
	return s.a
}
