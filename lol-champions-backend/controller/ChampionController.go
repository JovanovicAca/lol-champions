package controller

import (
	_ "database/sql"
	"encoding/json"
	"lol-champions-backend/helper"
	"lol-champions-backend/model"
	"lol-champions-backend/service"
	"net/http"
)

type ChampionController interface {
	GetAll(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
	DeleteChamp(response http.ResponseWriter, request *http.Request)
	UpdateChamp(response http.ResponseWriter, request *http.Request)
	FilterSearchChamps(response http.ResponseWriter, request *http.Request)
}
type championController struct {
}

var (
	championService service.ChampionService
)

func NewChampionController(service service.ChampionService) ChampionController {
	championService = service
	return &championController{}
}

func (*championController) FilterSearchChamps(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var searchFilter helper.FilterRequest

	err := json.NewDecoder(request.Body).Decode(&searchFilter)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := championService.SearchFilter(&searchFilter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (c *championController) UpdateChamp(response http.ResponseWriter, request *http.Request) {

	var champ model.Champion

	err := json.NewDecoder(request.Body).Decode(&champ)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}

	champion := model.Champion{
		Id:        champ.Id,
		Name:      champ.Name,
		Class:     champ.Class,
		World:     champ.World,
		Position:  champ.Position,
		Weapon:    champ.Weapon,
		MagicCost: champ.MagicCost,
	}
	response.Header().Set("Content-Type", "application/json")

	championService.UpdateChamp(champion)

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response)
}

func (*championController) DeleteChamp(response http.ResponseWriter, request *http.Request) {

	var champ model.Champion

	err := json.NewDecoder(request.Body).Decode(&champ)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}

	champion := model.Champion{
		Id:        champ.Id,
		Name:      champ.Name,
		Class:     champ.Class,
		World:     champ.World,
		Position:  champ.Position,
		Weapon:    champ.Weapon,
		MagicCost: champ.MagicCost,
	}
	response.Header().Set("Content-Type", "application/json")

	championService.DeleteChamp(champion)

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response)

}

func (*championController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var champs []model.Champion
	champs, _ = championService.GetAll()
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(champs)
}

func (*championController) Save(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var champ model.Champion

	err := json.NewDecoder(request.Body).Decode(&champ)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err1 := championService.Save(&champ)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}
