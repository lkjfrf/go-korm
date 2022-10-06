package korm

import (
	"database/sql"
	"log"
	"testing"
)

type Test_CreateTable struct {
	Id          int32  `korm:"integer"`
	Title       string `korm:"varchar(32)"`
	Description string `korm:"varchar(100)"`
}

func GetConnection(username string, password string, url string, schema string) *sql.DB {
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+url+")/"+schema)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

func TestCreateTable(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()
	model.CreateTable(db, true)
}

func TestGet(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()
	model.Data.Id = 777

	model.Get(db, 0)
	log.Println(model.Data.Id, model.Data.Title, model.Data.Description)
}

func TestGetAll(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()

	result := model.GetAll(db)
	log.Println(result)
}

func TestInsert(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()
	model.Data.Id = 777
	model.Data.Title = "TestInsertTitle777"
	model.Data.Description = "TestDesc777"

	model.Insert(db)
}

func TestUpdate(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()
	model.Data.Id = 777
	model.Data.Title = "*UpdatedTitle777"
	model.Data.Description = "*UpdatedDesc777"

	model.Update(db, 0)
}

func TestDelete(t *testing.T) {

	db := GetConnection("root", "password", "localhost:3306", "korm_test")

	model := NewModel[Test_CreateTable]()
	model.Data.Id = 777

	model.Delete(db, 0)
}
