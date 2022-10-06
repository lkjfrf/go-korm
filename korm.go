package korm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Model[T any] struct {
	tablename string
	fieldnum  int
	fieldname []string
	datatype  []string

	Data T
}

func NewModel[T any]() Model[T] {

	model := Model[T]{}

	val := reflect.Indirect(reflect.ValueOf(model.Data))
	tablename := val.Type().Name()

	model.tablename = tablename
	model.fieldnum = val.NumField()

	for i := 0; i < val.NumField(); i++ {
		fieldname := val.Type().Field(i).Name
		tagtype := val.Type().Field(i).Tag.Get("korm")

		model.fieldname = append(model.fieldname, fieldname)
		model.datatype = append(model.datatype, tagtype)
	}

	return model
}

func (m *Model[T]) CreateTable(db *sql.DB, withpk bool) {

	query := "CREATE TABLE if not exists " + m.tablename + " ("

	for i := 0; i < m.fieldnum; i++ {
		if i == 0 && withpk {
			query += m.fieldname[i] + " " + m.datatype[i] + " PRIMARY KEY" + ","
		} else {
			query += m.fieldname[i] + " " + m.datatype[i] + ","
		}
	}

	query = query[:len(query)-1]
	query += ")"

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func (m *Model[T]) Get(db *sql.DB, fieldidx int) bool {

	val := reflect.Indirect(reflect.ValueOf(m.Data))
	value := val.Field(fieldidx).Interface()

	query := "SELECT * FROM " + m.tablename + " WHERE " + m.fieldname[fieldidx] + "="

	strvalue := fmt.Sprintf("%v", value)
	if strings.Contains(m.datatype[fieldidx], "varchar") {
		query += "'" + strvalue + "'"
	} else {
		query += strvalue
	}

	result, _ := db.Query(query)
	jsons := QueryRowToJson(result)

	var datas []T
	err := json.Unmarshal([]byte(jsons), &datas)
	if err != nil {
		log.Println("GetByPK : json.Unmarshall error", err)
	}

	if len(datas) > 0 {
		m.Data = datas[0]
		return true
	}

	return false
}

func (m *Model[T]) GetAll(db *sql.DB) (results []T) {

	query := "SELECT * FROM " + m.tablename
	result, _ := db.Query(query)
	jsons := QueryRowToJson(result)

	var datas []T
	err := json.Unmarshal([]byte(jsons), &datas)
	if err != nil {
		log.Println("GetByPK : json.Unmarshall error", err)
	}

	results = append(results, datas...)
	return results
}

func (m *Model[T]) Insert(db *sql.DB) bool {

	val := reflect.Indirect(reflect.ValueOf(m.Data))

	query := "INSERT INTO " + m.tablename + " VALUES("
	for i := 0; i < m.fieldnum; i++ {

		value := val.Field(i).Interface()
		strvalue := fmt.Sprintf("%v", value)

		if strings.Contains(m.datatype[i], "varchar") {
			query += "'" + strvalue + "'"
		} else {
			query += strvalue
		}
		query += ","
	}

	query = query[:len(query)-1]
	query += ")"

	_, err := db.Exec(query)
	if err != nil {
		log.Println("Insert error", err)
		return false
	}

	return true
}

func (m *Model[T]) Update(db *sql.DB, fieldidx int) bool {

	val := reflect.Indirect(reflect.ValueOf(m.Data))

	query := "UPDATE " + m.tablename + " SET "
	for i := 0; i < m.fieldnum; i++ {
		value := val.Field(i).Interface()
		strvalue := fmt.Sprintf("%v", value)

		if strings.Contains(m.datatype[i], "varchar") {
			query += m.fieldname[i] + "=" + "'" + strvalue + "'"
		} else {
			query += m.fieldname[i] + "=" + strvalue
		}
		query += ","
	}

	query = query[:len(query)-1]
	query += " WHERE " + m.fieldname[fieldidx] + "="

	strvalue := fmt.Sprintf("%v", val.Field(fieldidx).Interface())
	if strings.Contains(m.datatype[fieldidx], "varchar") {
		query += "'" + strvalue + "'"
	} else {
		query += strvalue
	}

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (m *Model[T]) Delete(db *sql.DB, fieldidx int) bool {
	val := reflect.Indirect(reflect.ValueOf(m.Data))
	value := val.Field(fieldidx).Interface()
	strvalue := fmt.Sprintf("%v", value)

	query := "DELETE FROM " + m.tablename + " WHERE " + m.fieldname[fieldidx] + "="
	if strings.Contains(m.datatype[fieldidx], "varchar") {
		query += "'" + strvalue + "'"
	} else {
		query += strvalue
	}

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

/* ---------------
	Json Utils
--------------- */

func QueryRowToJson(rows *sql.Rows) string {
	columns, err := rows.ColumnTypes()
	if err != nil {
		panic(err.Error())
	}

	values := make([]any, len(columns))
	scanArgs := make([]any, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	c := 0
	results := make(map[string]any)
	data := []string{}

	for rows.Next() {
		if c > 0 {
			data = append(data, ",")
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for i, value := range values {

			if columns[i].DatabaseTypeName() == "VARCHAR" {

				s := string(value.([]byte))
				results[columns[i].Name()] = s

			} else if columns[i].DatabaseTypeName() == "INT" {

				s := string(value.([]byte))
				x, err := strconv.Atoi(s)
				if err == nil {
					results[columns[i].Name()] = x
				}

			} else {
				results[columns[i].Name()] = value
			}
		}

		b, err := json.Marshal(results)
		if err != nil {
			log.Println("QueryRowToJson : ERROR", err)
		}
		data = append(data, strings.TrimSpace(string(b)))
		c++
	}

	result := fmt.Sprint(data)
	return result
}
