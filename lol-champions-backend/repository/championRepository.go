package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"lol-champions-backend/helper"
	"lol-champions-backend/model"
	"strings"
)

type ChampionRepository interface {
	GetAll() ([]model.Champion, error)
	Save(champ model.Champion) (model.Champion, error)
	FindById(id uuid.UUID) model.Champion
	DeleteChamp(champ model.Champion)
	UpdateChamp(champ model.Champion)
	SearchFilter(champs []model.Champion, filter helper.FilterRequest) []model.Champion
	Filter(searched []model.Champion, filter helper.FilterRequest) []model.Champion
}

type championRepository struct {
}

func (r *championRepository) Filter(responseChamps []model.Champion, filter helper.FilterRequest) []model.Champion {
	//Filtering champions by class, positions, weapon, magic cost
	//var allChamps []model.Champion
	//allChamps, _ = r.GetAll()

	var returnedChamps []model.Champion
	var classChamps []model.Champion
	var weaponChamps []model.Champion
	var magicCostChamps []model.Champion
	var positionChamps []model.Champion

	for y, element1 := range responseChamps {
		if filter.Class != "" {
			if strings.Compare(element1.Class, filter.Class) == 0 {
				fmt.Println("A", responseChamps[y].Name)
				classChamps = append(classChamps, responseChamps[y])
			}
		}
		if filter.Weapon != "" {
			if strings.Compare(element1.Weapon, filter.Weapon) == 0 {
				fmt.Println("B", responseChamps[y].Name)
				weaponChamps = append(weaponChamps, responseChamps[y])
			}
		}
		if filter.MagicCost != "" {
			if strings.Compare(element1.MagicCost, filter.MagicCost) == 0 {
				fmt.Println("C", responseChamps[y].Name)
				magicCostChamps = append(magicCostChamps, responseChamps[y])

			}
		}
		if filter.Positions != nil {
			for _, el1 := range responseChamps[y].Position {
				for _, el := range filter.Positions {

					if strings.Compare(el1, el) == 0 {
						positionChamps = append(positionChamps, responseChamps[y])
					}
				}

			}
		}
	}
	returnedChamps = responseChamps
	if filter.Class != "" {
		returnedChamps = sameElements(returnedChamps, classChamps)
	}
	if filter.Weapon != "" {
		returnedChamps = sameElements(returnedChamps, weaponChamps)
	}
	if filter.MagicCost != "" {
		returnedChamps = sameElements(returnedChamps, magicCostChamps)
	}
	if filter.Positions != nil {
		returnedChamps = sameElements(returnedChamps, positionChamps)
	}
	return returnedChamps
}

func (*championRepository) SearchFilter(champs []model.Champion, filter helper.FilterRequest) []model.Champion {
	//filter -> object with search and filter attributes
	//champs -> all champions
	//respChamps -> champions that we will return

	fmt.Println("Searching and filtering")
	var respChamps []model.Champion
	respChamps = champs
	//First doing searching
	//Search by name
	var searchedChampsName []model.Champion
	if filter.NameSearch != "" {
		fmt.Println("name search")
		for i, element := range champs {
			if strings.Contains(element.Name, filter.NameSearch) {
				searchedChampsName = append(searchedChampsName, champs[i])
			}
		}
	}
	//Search by world
	var searchedChampsWorld []model.Champion
	if filter.WorldSearch != "" {
		fmt.Println("world search")
		for i, element := range champs {
			if strings.Contains(element.World, filter.WorldSearch) {
				searchedChampsWorld = append(searchedChampsWorld, champs[i])
			}
		}
	}

	//Make intersection between two lists
	if filter.NameSearch != "" && filter.WorldSearch != "" {
		respChamps = intersection(searchedChampsName, searchedChampsWorld)
	} else if filter.NameSearch != "" && filter.WorldSearch == "" {
		respChamps = searchedChampsName
	} else if filter.NameSearch == "" && filter.WorldSearch != "" {
		respChamps = searchedChampsWorld
	} else {
		respChamps = champs
	}

	return respChamps
}

func sameElements(champs []model.Champion, list []model.Champion) []model.Champion {
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
	return diff
}

func difference(champs []model.Champion, list []model.Champion) []model.Champion {
	mb := make(map[uuid.UUID]struct{}, len(list))
	for _, x := range list {
		mb[x.Id] = struct{}{}
	}
	var diff []model.Champion
	for i, x := range champs {
		if _, found := mb[x.Id]; !found {
			diff = append(diff, champs[i])
		}
	}
	return diff
}

func intersection(name []model.Champion, world []model.Champion) []model.Champion {
	out := []model.Champion{}
	bucket := map[uuid.UUID]bool{}

	for _, i := range name {
		for _, j := range world {
			if i.Id == j.Id && !bucket[i.Id] {
				out = append(out, i)
				bucket[i.Id] = true
			}
		}
	}
	return out
}

func (*championRepository) UpdateChamp(champ model.Champion) {
	fmt.Println("Updating champion")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	defer sqlObj.Close()

	updatedChamp := `UPDATE "Champion" SET "name" = $2, "world" = $3, "class" =  $4, "weapon" = $5, "MagicCost" = $6 WHERE "id" = $1`
	_, error := sqlObj.Exec(updatedChamp, champ.Id, champ.Name, champ.World, champ.Class, champ.Weapon, champ.MagicCost)
	if error != nil {
		panic(error)
	}

	//Change in Champion_position champs positions
	//First delete all positions and then add new
	for range champ.Position {
		deleted := `DELETE FROM "Champion_position" WHERE "championId" = $1`
		_, error := sqlObj.Exec(deleted, champ.Id)
		if error != nil {
			panic(error)
		}

	}
	for _, element := range champ.Position {
		position := NewPositionRepository().FindByName(element)
		insert := `insert into "Champion_position"("championId","positionId") values ($1, $2)`
		_, e := sqlObj.Exec(insert, champ.Id, position.Id)
		if e != nil {
			panic(e)
		}
	}
	fmt.Println("Successfully Updated Champion!")
}

func (r *championRepository) DeleteChamp(champ model.Champion) {
	fmt.Println("Deleting Champion")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	defer sqlObj.Close()

	//Delete from Champion_position also
	deletedChampPos := `DELETE FROM "Champion_position" WHERE "championId" = $1`
	_, error1 := sqlObj.Exec(deletedChampPos, champ.Id)
	if error1 != nil {
		panic(error1)
	}
	deletedChamp := `DELETE FROM "Champion" WHERE "id" = $1`
	res, error := sqlObj.Exec(deletedChamp, champ.Id)
	if error != nil {
		panic(error)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
	fmt.Println("Successfully Deleted Champion!")
}

func (r *championRepository) FindById(id uuid.UUID) model.Champion {
	//Champ for return
	champ := model.Champion{}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	//Find champion
	selection, error := sqlObj.Query(`SELECT * FROM "Champion" WHERE "Id" = $1`, id)

	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var name string
		var world string
		var class string
		var weapon string
		var magicCost string

		error1 := selection.Scan(&id, &name, &world, &class, &weapon, &magicCost)
		if error1 != nil {
			panic(error1)
		}

		champ.Id = id
		champ.Name = name
		champ.World = world
		champ.Class = class
		champ.Weapon = weapon
		champ.MagicCost = magicCost

		//Find champions positions
		positionsSelection, error2 := sqlObj.Query(`SELECT * FROM "Champion_position" WHERE "championId" = $1`, id)

		if error2 != nil {
			panic(error2)
		}

		//List for champion positions
		var positionsChamp []model.Position
		var positionsChampString []string

		//From champion_positions take positions from Positions
		for positionsSelection.Next() {
			var pId uuid.UUID
			var positionName string

			error3 := positionsSelection.Scan(&pId, &positionName)
			if error3 != nil {
				panic(error3)
			}

			//Find positions
			positionSelection, error4 := sqlObj.Query(`SELECT * FROM "Position" where "id" = $1`, pId)

			if error4 != nil {
				panic(error4)
			}

			for positionSelection.Next() {
				var poId uuid.UUID
				var poName string
				error5 := positionSelection.Scan(&poId, &poName)
				if error5 != nil {
					panic(error5)
				}
				position := model.Position{
					Id:       poId,
					Position: poName,
				}
				//Adding positions to the list
				positionsChamp = append(positionsChamp, position)
			}
			//Add positions as strings in champ
			for _, element := range positionsChamp {
				positionsChampString = append(positionsChampString, element.Position)
			}
			champ.Position = positionsChampString
		}
	}
	return champ
}

func (*championRepository) GetAll() ([]model.Champion, error) {
	var champList []model.Champion

	fmt.Println("Getting All Champions")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
		panic(connectionError)
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "Champion"`)
	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var name string
		var world string
		var class string
		var weapon string
		var magicCost string

		error1 := selection.Scan(&id, &name, &world, &class, &weapon, &magicCost)
		if error1 != nil {
			panic(error1)
		}
		//Find all positions id from Champion_position
		var positionsChampString []string
		var positionsChamp []model.Position
		positionsIds := FindPositionIdsFromChampionId(id)

		for _, element := range positionsIds {
			positionsChamp = append(positionsChamp, PositionRepository(NewPositionRepository()).FindById(element))
		}
		for _, element1 := range positionsChamp {
			positionsChampString = append(positionsChampString, element1.Position)
		}
		champ := model.Champion{
			Id:        id,
			Name:      name,
			World:     world,
			Class:     class,
			Position:  positionsChampString,
			Weapon:    weapon,
			MagicCost: magicCost,
		}
		champList = append(champList, champ)
		//champ.Position = positionsChampString
	}
	return champList, nil
}

func FindPositionIdsFromChampionId(id uuid.UUID) []uuid.UUID {
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	selection, error := sqlObj.Query(`SELECT * FROM "Champion_position" WHERE "championId" = $1`, id)

	if error != nil {
		panic(error)
	}

	var idList []uuid.UUID

	for selection.Next() {
		var cId uuid.UUID
		var pId uuid.UUID
		error1 := selection.Scan(&cId, &pId)
		if error1 != nil {
			panic(error1)
		}
		idList = append(idList, pId)
	}

	return idList
}

func (*championRepository) Save(champ model.Champion) (model.Champion, error) {
	fmt.Println("Adding New Champion")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	// Insert into db
	var id = champ.Id.String()

	insert := `insert into "Champion"("id", "name", "world", "class", "weapon", "MagicCost") values ($1, $2, $3, $4, $5, $6)`
	//fmt.Println(champ)
	_, e := sqlObj.Exec(insert, id, champ.Name, champ.World, champ.Class, champ.Weapon, champ.MagicCost)
	if e != nil {
		fmt.Println(e)
		return champ, e
	}

	//Add champs positions
	for _, element := range champ.Position {
		position := NewPositionRepository().FindByName(element)
		insert := `insert into "Champion_position"("championId","positionId") values ($1, $2)`
		_, e := sqlObj.Exec(insert, champ.Id, position.Id)
		if e != nil {
			return champ, e
		}
	}
	defer sqlObj.Close()
	fmt.Println("Successfully added champion!")
	return champ, nil
}

func NewChampionRepository() ChampionRepository {
	return &championRepository{}
}
