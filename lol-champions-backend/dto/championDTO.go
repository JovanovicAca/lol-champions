package dto

import (
	"github.com/google/uuid"
)

type ChampionDTO struct {
	Id        uuid.UUID
	Name      string
	World     uuid.UUID
	Class     string
	Position  []string
	Weapon    string
	MagicCost string
}
