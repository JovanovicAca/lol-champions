package repository

import (
	"database/sql"
	"fmt"
	"lol-champions-backend/model"

	"github.com/google/uuid"
)

type PositionRepository interface {
	Save(position model.Position) (model.Position, error)
	FindByName(name string) model.Position
	FindById(id uuid.UUID) model.Position
}

type positionRepository struct {
}

func NewPositionRepository() PositionRepository {
	return &positionRepository{}
}

func (r *positionRepository) FindByName(name string) model.Position {
	position := model.Position{}

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "positions" WHERE "position" = $1`, name)

	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var p string
		error1 := selection.Scan(&id, &p)
		if error1 != nil {
			panic(error1)
		}
		position.Id = id
		position.Position = p
	}
	return position
}

func (*positionRepository) FindById(id uuid.UUID) model.Position {
	position := model.Position{}

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "positions" WHERE "id" = $1`, id)

	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var p string
		error1 := selection.Scan(&id, &p)
		if error1 != nil {
			panic(error1)
		}
		position.Id = id
		position.Position = p
	}
	return position
}

func (*positionRepository) Save(position model.Position) (model.Position, error) {
	fmt.Println("Adding New Position")

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	var id = position.Id.String()

	insert := `insert into "positions"("id","position") values ($1, $2)`
	_, e := sqlObj.Exec(insert, id, position.Position)
	if e != nil {
		return position, e
	}
	defer sqlObj.Close()
	fmt.Println("Successfully added position")
	return position, nil
}
