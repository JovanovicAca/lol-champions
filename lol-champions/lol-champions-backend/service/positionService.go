package service

import (
	"github.com/google/uuid"
	"lol-champions-backend/model"
	"lol-champions-backend/repository"
)

type PositionService interface {
	Save(position *model.Position) (*model.Position, error)
}

type positionService struct {
}

var (
	positionRepository repository.PositionRepository
)

func NewPositionService(positionRepo repository.PositionRepository) PositionService {
	positionRepository = positionRepo
	return &positionService{}
}

func (*positionService) Save(position *model.Position) (*model.Position, error) {
	p := model.Position{
		Id:       uuid.New(),
		Position: position.Position,
	}
	resp, err := positionRepository.Save(p)
	return &resp, err
}
