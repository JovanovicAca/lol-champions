package model

import (
	"github.com/google/uuid"
)

type Champion struct {
	Id        uuid.UUID
	Name      string
	World     string
	Class     string
	Position  []string
	Weapon    string
	MagicCost string
}
