package service

import (
	"github.com/google/uuid"
	"lol-champions-backend/helper"
	"lol-champions-backend/model"
	"lol-champions-backend/repository"
)

type ChampionService interface {
	GetAll() ([]model.Champion, error)
	Save(champ *model.Champion) (*model.Champion, error)
	DeleteChamp(champ model.Champion)
	UpdateChamp(champ model.Champion)
	SearchFilter(searchFilter *helper.FilterRequest) (*[]model.Champion, error)
}

type championService struct {
}

var (
	championRepository repository.ChampionRepository
)

func NewChampService(championRepo repository.ChampionRepository) ChampionService {
	championRepository = championRepo
	return &championService{}
}

func (*championService) SearchFilter(searchFilter *helper.FilterRequest) (*[]model.Champion, error) {
	var responseChamps []model.Champion
	//Get all and search/filter in that list
	responseChamps, _ = championRepository.GetAll()
	//List where everything will be stored
	searched := championRepository.SearchFilter(responseChamps, *searchFilter)
	filtered := championRepository.Filter(responseChamps, *searchFilter)
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

func (s *championService) UpdateChamp(champ model.Champion) {
	championRepository.UpdateChamp(champ)
}

func (*championService) DeleteChamp(champ model.Champion) {
	championRepository.DeleteChamp(champ)
}

func (*championService) Save(champ *model.Champion) (*model.Champion, error) {
	champion := model.Champion{
		Id:        uuid.New(),
		Name:      champ.Name,
		World:     champ.World,
		Class:     champ.Class,
		Position:  champ.Position,
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
