package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"lol-champions-backend/model"
)

type WorldRepository interface {
	FindAll() ([]model.World, error)
	Save(world model.World) (model.World, error)
}

type worldRepository struct {
}

func NewWorldRepository() WorldRepository {
	return &worldRepository{}
}

func (*worldRepository) FindAll() ([]model.World, error) {
	//TODO implement me
	panic("implement me")
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

	insert := `insert into "World"("id","name","description") values ($1, $2,$3)`
	_, e := sqlObj.Exec(insert, id, world.Name, world.Description)
	if e != nil {
		return world, e
	}

	defer sqlObj.Close()
	fmt.Println("Successfully added world!")
	return world, nil
}
