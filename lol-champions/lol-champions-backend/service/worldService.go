package service

import (
	"github.com/google/uuid"
	"lol-champions-backend/model"
	"lol-champions-backend/repository"
)

type WorldService interface {
	GetAll([]model.World, error)
	Save(world *model.World) (*model.World, error)
}

type worldService struct {
}

var (
	worldRepository repository.WorldRepository
)

func NewWorldService(worldRepo repository.WorldRepository) WorldService {
	worldRepository = worldRepo
	return &worldService{}
}

func (*worldService) GetAll(worlds []model.World, err error) {
	//TODO implement me
	panic("implement me")
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
