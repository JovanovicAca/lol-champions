package main

import (
	_ "bufio"
	_ "database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"lol-champions-backend/controller"
	"lol-champions-backend/repository"
	"lol-champions-backend/router"
	"lol-champions-backend/service"
	"net/http"
)

func main() {
	var (
		championRepository repository.ChampionRepository = repository.NewChampionRepository()
		worldRepository    repository.WorldRepository    = repository.NewWorldRepository()
		positionRepository repository.PositionRepository = repository.NewPositionRepository()

		championService service.ChampionService = service.NewChampService(championRepository)
		worldService    service.WorldService    = service.NewWorldService(worldRepository)
		positionService service.PositionService = service.NewPositionService(positionRepository)

		championController controller.ChampionController = controller.NewChampionController(championService)
		worldController    controller.WorldController    = controller.NewWorldController(worldService)
		positionController controller.PositionController = controller.NewPositionController(positionService)
	)

	var httpRouter router.Router = router.NewMuxRouter()
	//http.ListenAndServe(":8080", nil)
	fmt.Println("Hello world!")
	const port string = ":8080"
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Up and runing")
	})

	httpRouter.POST("/addChampion", championController.Save)
	httpRouter.GET("/getAll", championController.GetAll)
	httpRouter.POST("/deleteChamp", championController.DeleteChamp)
	httpRouter.POST("/updateChamp", championController.UpdateChamp)
	httpRouter.POST("/searchFilter", championController.FilterSearchChamps)
	httpRouter.POST("/addWorld", worldController.Save)
	httpRouter.POST("/addPosition", positionController.Save)
	httpRouter.SERVE(port)
}
