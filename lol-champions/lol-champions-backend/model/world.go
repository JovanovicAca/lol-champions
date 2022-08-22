package model

import (
	"github.com/google/uuid"
)

type World struct {
	Id          uuid.UUID
	Name        string
	Description string
}
