package controller

import (
	"encoding/json"
	"lol-champions-backend/model"
	"lol-champions-backend/service"
	"net/http"
)

type PositionController interface {
	Save(response http.ResponseWriter, request *http.Request)
}

type positionController struct {
}

var (
	positionService service.PositionService
)

func NewPositionController(service service.PositionService) PositionController {
	positionService = service
	return &positionController{}
}

func (*positionController) Save(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var position model.Position

	err := json.NewDecoder(request.Body).Decode(&position)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err1 := positionService.Save(&position)
	if err1 != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(result)
}
