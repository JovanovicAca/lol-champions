package service

import (
	"lol-champions-backend/model"
	"lol-champions-backend/repository"

	"github.com/google/uuid"
)

type WorldService interface {
	GetAll() ([]model.World, error)
	Save(world *model.World) (*model.World, error)
}

type worldService struct {
}

var (
//worldRepository repository.WorldRepository
)

func NewWorldService(worldRepo repository.WorldRepository) WorldService {
	worldRepository = worldRepo
	return &worldService{}
}

func (w *worldService) GetAll() (worlds []model.World, err error) {
	var response []model.World
	response, _ = worldRepository.GetAll()
	return response, nil
}

func (*worldService) Save(world *model.World) (*model.World, error) {
	w := model.World{
		Id:          uuid.New(),
		Name:        world.Name,
		Description: world.Description,
	}
	resp, err := worldRepository.Save(w)
	return &resp, err
}
