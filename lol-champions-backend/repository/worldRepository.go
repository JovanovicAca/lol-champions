package repository

import (
	"database/sql"
	"fmt"
	"lol-champions-backend/model"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type WorldRepository interface {
	GetAll() ([]model.World, error)
	Save(world model.World) (model.World, error)
	FindById(id uuid.UUID) model.World
	FindByName(name string) model.World
}

type worldRepository struct {
}

func NewWorldRepository() WorldRepository {
	return &worldRepository{}
}

func (w *worldRepository) FindByName(name string) model.World {
	world := model.World{}

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "worlds" WHERE "name" = $1`, name)
	if error != nil {
		panic(error)
	}
	for selection.Next() {
		var id uuid.UUID
		var name string
		var description string
		error1 := selection.Scan(&id, &name, &description)
		if error1 != nil {
			panic(error1)
		}
		world.Id = id
		world.Name = name
		world.Description = description
	}
	return world
}

func (*worldRepository) GetAll() ([]model.World, error) {
	var worldList []model.World
	fmt.Println("Getting All Worlds")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		panic(connectionError)
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "worlds"`)
	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var name string
		var description string

		error1 := selection.Scan(&id, &name, &description)
		if error1 != nil {
			panic(error1)
		}

		world := model.World{
			Id:          id,
			Name:        name,
			Description: description,
		}
		worldList = append(worldList, world)
	}
	return worldList, nil
}

func (*worldRepository) Save(world model.World) (model.World, error) {
	fmt.Println("Adding New World")

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	var id = world.Id.String()

	insert := `insert into "worlds"("id","name","description") values ($1, $2,$3)`
	_, e := sqlObj.Exec(insert, id, world.Name, world.Description)
	if e != nil {
		return world, e
	}

	defer sqlObj.Close()
	fmt.Println("Successfully added world!")
	return world, nil
}

func (r *worldRepository) FindById(id uuid.UUID) model.World {
	world := model.World{}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)
	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	//Find world
	selection, error := sqlObj.Query(`SELECT * FROM "worlds" WHERE "id" = $1`, id)
	if error != nil {
		panic(error)
	}
	if selection.Next() {
		var id uuid.UUID
		var name string
		var description string
		error1 := selection.Scan(&id, &name, &description)
		if error1 != nil {
			panic(error1)
		}
		world.Id = id
		world.Name = name
		world.Description = description
	}
	return world
}
