package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Champion struct {
	Id        uuid.UUID
	Name      string
	World     World
	Class     string
	Position  []Position
	Weapon    string
	MagicCost string
}

func (c Champion) String() string {
	return fmt.Sprintf("%v %v %v %v %v %v %v", c.Id, c.Name, c.World.Id, c.Class, c.Position, c.Weapon, c.MagicCost)
}
