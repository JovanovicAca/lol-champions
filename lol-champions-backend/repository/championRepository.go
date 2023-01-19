package repository

import (
	"database/sql"
	"fmt"
	"lol-champions-backend/helper"
	"lol-champions-backend/model"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ChampionRepository interface {
	GetAll() ([]model.Champion, error)
	Save(champ model.Champion) (model.Champion, error)
	FindById(id uuid.UUID) (model.Champion, error)
	DeleteChamp(id uuid.UUID) int
	UpdateChamp(champ model.Champion) (model.Champion, error)
	SearchFilter(champs []model.Champion, filter helper.FilterRequest) []model.Champion
	Filter(searched []model.Champion, filter helper.FilterRequest) []model.Champion
}

type championRepository struct {
}

func NewChampionRepository() ChampionRepository {
	return &championRepository{}
}

func (r *championRepository) Filter(champs []model.Champion, filter helper.FilterRequest) []model.Champion {
	//Filtering champions by class, positions, weapon, magic cost

	var returnedChamps []model.Champion
	var classChamps []model.Champion
	var weaponChamps []model.Champion
	var magicCostChamps []model.Champion
	var positionChamps []model.Champion

	for y, element := range champs {
		if filter.Class != "" {
			if strings.Compare(element.Class, filter.Class) == 0 {
				classChamps = append(classChamps, champs[y])
			}
		}
		if filter.Weapon != "" {
			if strings.Compare(element.Weapon, filter.Weapon) == 0 {
				weaponChamps = append(weaponChamps, champs[y])
			}
		}
		if filter.MagicCost != "" {
			if strings.Compare(element.MagicCost, filter.MagicCost) == 0 {
				magicCostChamps = append(magicCostChamps, champs[y])

			}
		}
		if filter.Positions != nil {
			//List of champs positions needed for position filtering
			var champsPositions = element.Position
			for _, el := range filter.Positions {
				if contains(champsPositions, el) {
					positionChamps = append(positionChamps, champs[y])
				}
			}
		}
	}
	returnedChamps = champs
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

func contains(champsPositions []model.Position, el string) bool {
	for _, e := range champsPositions {
		if e.Position == el {
			return true
		}
	}
	return false
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
		fmt.Println("Name search")
		for i, element := range champs {
			if strings.Contains(element.Name, filter.NameSearch) {
				searchedChampsName = append(searchedChampsName, champs[i])
			}
		}
	}
	//Search by world
	var searchedChampsWorld []model.Champion
	if filter.WorldSearch != "" {
		fmt.Println("World search")
		for i, element := range champs {
			if strings.Contains(element.World.Name, filter.WorldSearch) {
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

// func difference(champs []model.Champion, list []model.Champion) []model.Champion {
// 	mb := make(map[uuid.UUID]struct{}, len(list))
// 	for _, x := range list {
// 		mb[x.Id] = struct{}{}
// 	}
// 	var diff []model.Champion
// 	for i, x := range champs {
// 		if _, found := mb[x.Id]; !found {
// 			diff = append(diff, champs[i])
// 		}
// 	}
// 	return diff
// }

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

func (*championRepository) UpdateChamp(champ model.Champion) (model.Champion, error) {
	fmt.Println("Updating champion")

	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	var worldId = champ.World.Id.String()
	updatedChamp := `UPDATE "champions" SET "name" = $2, "class" =  $3, "weapon" = $4, "magiccost" = $5, "worldsid" = $6 WHERE "id" = $1`
	_, error := sqlObj.Exec(updatedChamp, champ.Id, champ.Name, champ.Class, champ.Weapon, champ.MagicCost, worldId)
	if error != nil {
		fmt.Println("2")
		panic(error)

	}

	//Change in champion_position champs positions
	//First delete all positions and then add new
	for range champ.Position {
		deleted := `DELETE FROM "champion_position" WHERE "championid" = $1`
		_, error := sqlObj.Exec(deleted, champ.Id)
		if error != nil {
			panic(error)
		}
	}
	for _, element := range champ.Position {
		position := NewPositionRepository().FindByName(element.Position)
		insert := `insert into "champion_position"("championid","positionid") values ($1, $2)`
		_, e := sqlObj.Exec(insert, champ.Id, position.Id)
		if e != nil {
			return champ, e
		}
	}
	defer sqlObj.Close()
	fmt.Println("Successfully Updated Champion!")
	return champ, nil
}

func (r *championRepository) DeleteChamp(id uuid.UUID) int {
	fmt.Println("Deleting Champion")
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	defer sqlObj.Close()

	//Delete from Champion_position also
	deletedChampPos := `DELETE FROM "champion_position" WHERE "championid" = $1`
	_, error1 := sqlObj.Exec(deletedChampPos, id)
	if error1 != nil {
		return 1
	}
	deletedChamp := `DELETE FROM "champions" WHERE "id" = $1`
	res, error := sqlObj.Exec(deletedChamp, id)
	if error != nil {
		return 1
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 1
	}
	if count == 0 {
		return 1
	}
	fmt.Println("Successfully Deleted Champion!")
	return 0
}

func (r *championRepository) FindById(id uuid.UUID) (model.Champion, error) {
	//Champ for return
	//Getting champion by id
	champ := model.Champion{}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}

	defer sqlObj.Close()

	//Find champion
	selection, error := sqlObj.Query(`SELECT * FROM "champions" WHERE "id" = $1`, id)

	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var name string
		var class string
		var weapon string
		var magicCost string
		var worldsid uuid.UUID

		error1 := selection.Scan(&id, &name, &class, &weapon, &magicCost, &worldsid)
		if error1 != nil {
			panic(error1)
		}

		var positionsChamp []model.Position
		positionsIds := FindPositionIdsFromChampionId(id)
		for _, element := range positionsIds {
			positionsChamp = append(positionsChamp, NewPositionRepository().FindById(element))
		}
		var world = WorldRepository.FindById(NewWorldRepository(), worldsid)
		champ.Id = id
		champ.Name = name
		champ.World = world
		champ.Class = class
		champ.Position = positionsChamp
		champ.Weapon = weapon
		champ.MagicCost = magicCost

	}
	return champ, nil
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

	selection, error := sqlObj.Query(`SELECT * FROM "champions"`)
	if error != nil {
		panic(error)
	}

	for selection.Next() {
		var id uuid.UUID
		var name string
		var class string
		var weapon string
		var magicCost string
		var worldsid uuid.UUID

		error1 := selection.Scan(&id, &name, &class, &weapon, &magicCost, &worldsid)
		if error1 != nil {
			panic(error1)
		}
		//Find all positions id from champion_position
		var positionsChamp []model.Position
		positionsIds := FindPositionIdsFromChampionId(id)

		for _, element := range positionsIds {
			positionsChamp = append(positionsChamp, NewPositionRepository().FindById(element))
		}

		var world = WorldRepository.FindById(NewWorldRepository(), worldsid)
		champ := model.Champion{
			Id:        id,
			Name:      name,
			World:     world,
			Class:     class,
			Position:  positionsChamp,
			Weapon:    weapon,
			MagicCost: magicCost,
		}
		champList = append(champList, champ)
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

	selection, error := sqlObj.Query(`SELECT * FROM "champion_position" WHERE "championid" = $1`, id)

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
	fmt.Println(champ)
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	sqlObj, connectionError := sql.Open("postgres", sqlConn)

	if connectionError != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
	}
	// Insert into db
	var id = champ.Id.String()
	var worldId = champ.World.Id.String()
	insert := `insert into "champions"("id", "name", "class", "weapon", "magiccost","worldsid") values ($1, $2, $3, $4, $5, $6)`
	_, e := sqlObj.Exec(insert, id, champ.Name, champ.Class, champ.Weapon, champ.MagicCost, worldId)
	if e != nil {
		fmt.Println(e)
		return champ, e
	}

	//Add champs positions
	for _, element := range champ.Position {
		position := NewPositionRepository().FindByName(element.Position)
		insert := `insert into "champion_position"("championid","positionid") values ($1, $2)`
		_, e := sqlObj.Exec(insert, champ.Id, position.Id)
		if e != nil {
			return champ, e
		}
	}
	defer sqlObj.Close()
	fmt.Println("Successfully added champion!")
	return champ, nil
}
