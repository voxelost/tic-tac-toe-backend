package utils

import "github.com/google/uuid"

type ID string

func NewId() *ID {
	id := ID(uuid.NewString())
	return &id
}

func (id *ID) GetId() ID {
	return *id
}
