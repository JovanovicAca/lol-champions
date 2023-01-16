package controller

import (
	"encoding/json"
	"lol-champions-backend/model"
	"lol-champions-backend/service"
	"net/http"
)

type WorldController interface {
	GetAll(response http.ResponseWriter, request *http.Request)
	Save(response http.ResponseWriter, request *http.Request)
	//	GetById(response http.ResponseWriter, request *http.Request)
}

type worldController struct {
}

var (
	worldService service.WorldService
)

func NewWorldController(service service.WorldService) WorldController {
	worldService = service
	return &worldController{}
}

func (*worldController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var worlds []model.World
	worlds, _ = worldService.GetAll()
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(worlds)
}

func (*worldController) Save(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var world model.World

	err := json.NewDecoder(request.Body).Decode(&world)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err1 := worldService.Save(&world)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}

// func (*worldController) GetById(response http.ResponseWriter, request *http.Request) {
// 	//TODO implement me
// 	panic("implement me")
// }
