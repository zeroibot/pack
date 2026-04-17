package model

import "github.com/zeroibot/pack/clock"

type (
	ID       = uint
	Date     = clock.Date
	DateTime = clock.DateTime
)

// IDField is an embeddable ID property
type IDField struct {
	ID ID `json:"-"`
}

func (x IDField) GetID() ID {
	return x.ID
}

// CodeField is an embeddable Code property
type CodeField struct {
	Code string
}

func (x CodeField) GetCode() string {
	return x.Code
}

// CreatedAtField is an embeddable CreatedAt property
type CreatedAtField struct {
	CreatedAt DateTime
}

func (x CreatedAtField) GetDateTime() DateTime {
	return x.CreatedAt
}

// IsActiveField is an embeddable IsActive property
type IsActiveField struct {
	IsActive bool `json:"-"`
}

func (x IsActiveField) GetIsActive() bool {
	return x.IsActive
}

// Identity combines IDField and CodeField
type Identity struct {
	IDField
	CodeField
}

// AutoItem combines IDField, CreatedAtField, IsActiveField
type AutoItem struct {
	IDField
	CreatedAtField
	IsActiveField
}

// Item combines IDField, CreatedAtField, IsActiveField, CodeField
type Item struct {
	AutoItem
	CodeField
}

// Initialize sets the ID, CreatedAt, IsActive to default values
func (x *Item) Initialize() {
	x.AutoItem.Initialize()
}

// Initialize sets the ID, CreatedAt, IsActive to default values
func (x *AutoItem) Initialize() {
	x.ID = 0 // for auto-increment
	x.CreatedAt = clock.DateTimeNow()
	x.IsActive = true
}
