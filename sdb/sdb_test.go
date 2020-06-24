package sdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var stestdir string = "../../erp/data/Db/"
var stpath string = stestdir + "main.Db"

func TempFilename(t *testing.T) string {
	f, err := ioutil.TempFile("", "sdb-test-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func Test_GetTableNames(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int, y float)")
	tables, err := GetTableNames(db.Db)
	fmt.Println("tables:", tables)
}
func TestSdb_QueryValue(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int, y float)")
	cols, err := db.QueryValue("SELECT * FROM user")
	if err != nil {
		t.Fatal("Failed to GetColumns:", err)
	}
	fmt.Println("cols, ", cols)
}
func Test_GetColumns(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int not null , y float)")
	cols, err := GetColumns(db.Db, "user")
	if err != nil {
		t.Fatal("Failed to GetColumns:", err)
	}
	fmt.Println("cols, ", cols)
}
func TestSdb_GetTableInfos(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int, y float)")
	tableInfos, err := db.GetTableInfoDict()
	if err != nil {
		t.Fatal("Failed to GetTableInfoDict:", err)
	}
	fmt.Println("tableInfos:", tableInfos)
}
func TestDbase_AddSpecialCols(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int, y float)")
	if err := db.AddSpecialCols(); err != nil {
		t.Fatal("Failed to TestDbase_AddSpecialCols:", err)
	}
	fmt.Println("TestDbase_AddSpecialCols, ok;")
}

func TestDbase_RealAddSpecialCols(t *testing.T) {
	db, err := OpenDb("E:/gopath/src/etnet/data/db/main.db")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	if err := db.AddSpecialCols(); err != nil {
		t.Fatal("Failed to TestDbase_AddSpecialCols:", err)
	}
	fmt.Println("TestDbase_AddSpecialCols, ok;")
}
