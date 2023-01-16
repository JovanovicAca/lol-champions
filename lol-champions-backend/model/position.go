package model

import (
	"fmt"

	"github.com/google/uuid"
)

type Position struct {
	Id       uuid.UUID
	Position string
}

func (p Position) String() string {
	return fmt.Sprintf("%v %v", p.Id, p.Position)
}
