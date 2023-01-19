package service

import (
	"lol-champions-backend/dto"
	"lol-champions-backend/helper"
	"lol-champions-backend/model"
	"lol-champions-backend/repository"

	"github.com/google/uuid"
)

type ChampionService interface {
	GetAll() ([]model.Champion, error)
	Save(champ *dto.ChampionDTO) (*model.Champion, error)
	DeleteChamp(champ *dto.ChampionDTO) int
	UpdateChamp(champ *dto.ChampionDTO) (*model.Champion, error)
	SearchFilter(searchFilter *helper.FilterRequest) (*[]model.Champion, error)
}

type championService struct {
}

var (
	championRepository  repository.ChampionRepository
	worldRepository     repository.WorldRepository
	positionsRepository repository.PositionRepository
)

func NewChampService(championRepo repository.ChampionRepository, worldRepo repository.WorldRepository, posRepo repository.PositionRepository) ChampionService {
	championRepository = championRepo
	worldRepository = worldRepo
	positionRepository = posRepo
	return &championService{}
}

func (*championService) SearchFilter(searchFilter *helper.FilterRequest) (*[]model.Champion, error) {
	var responseChamps []model.Champion
	responseChamps, _ = championRepository.GetAll()
	//List where everything will be stored
	searched := championRepository.SearchFilter(responseChamps, *searchFilter)
	//Pass searched champs into filtering (if there is no searchs full list is passed)
	filtered := championRepository.Filter(searched, *searchFilter)
	return sameElements(searched, filtered)
}

func sameElements(champs []model.Champion, list []model.Champion) (*[]model.Champion, error) {
	mb := make(map[uuid.UUID]struct{}, len(list))
	for _, x := range list {
		mb[x.Id] = struct{}{}
	}
	var diff []model.Champion
	for i, x := range champs {
		if _, found := mb[x.Id]; found {
			diff = append(diff, champs[i])
		}
	}
	return &diff, nil
}

func (s *championService) UpdateChamp(champ *dto.ChampionDTO) (*model.Champion, error) {
	var world = worldRepository.FindById(champ.World)
	var positions []model.Position
	for _, p := range champ.Position {
		var position = positionRepository.FindByName(p)
		positions = append(positions, position)
	}
	champion := model.Champion{
		Id:        champ.Id,
		Name:      champ.Name,
		World:     world,
		Class:     champ.Class,
		Position:  positions,
		Weapon:    champ.Weapon,
		MagicCost: champ.MagicCost,
	}
	resp, err := championRepository.UpdateChamp(champion)
	return &resp, err
}

func (*championService) DeleteChamp(champ *dto.ChampionDTO) int {
	return championRepository.DeleteChamp(champ.Id)
}

func (*championService) Save(champ *dto.ChampionDTO) (*model.Champion, error) {
	var world = worldRepository.FindById(champ.World)
	var positions []model.Position
	for _, p := range champ.Position {
		var position = positionRepository.FindByName(p)
		positions = append(positions, position)
	}
	champion := model.Champion{
		Id:        uuid.New(),
		Name:      champ.Name,
		World:     world,
		Class:     champ.Class,
		Position:  positions,
		Weapon:    champ.Weapon,
		MagicCost: champ.MagicCost,
	}
	resp, err := championRepository.Save(champion)
	return &resp, err
}

func (s *championService) GetAll() (champions []model.Champion, err error) {
	var response []model.Champion
	response, _ = championRepository.GetAll()
	return response, nil
}
