package controller

import (
	_ "database/sql"
	"encoding/json"
	"lol-champions-backend/dto"
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
	//FindById(response http.ResponseWriter, request *http.Request)
	// FilterSearchChamps(response http.ResponseWriter, request *http.Request)
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

// func (*championController) FindById(response http.ResponseWriter, request *http.Request) {
// 	var id uuid.UUID
// 	fmt.Println("1")
// 	data := json.NewDecoder(request.Body)
// 	fmt.Println(data)
// 	err := data.Decode(&id)
// 	if err != nil {
// 		fmt.Println(err)
// 		response.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Println("2")

// 	// result, err1 := championService.FindById(&id)
// 	// if err1 != nil {
// 	// 	response.WriteHeader(http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	response.WriteHeader(http.StatusOK)
// 	json.NewEncoder(response)
// }
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

	response.Header().Set("Content-Type", "application/json")
	var champDTO dto.ChampionDTO
	err := json.NewDecoder(request.Body).Decode(&champDTO)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err1 := championService.UpdateChamp(&champDTO)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

func (*championController) DeleteChamp(response http.ResponseWriter, request *http.Request) {

	//var champ model.Champion
	var champ dto.ChampionDTO
	err := json.NewDecoder(request.Body).Decode(&champ)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
	}
	if championService.DeleteChamp(&champ) == 0 {
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response)
	} else {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response)
	}
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
	var champDTO dto.ChampionDTO
	err := json.NewDecoder(request.Body).Decode(&champDTO)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err1 := championService.Save(&champDTO)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}
