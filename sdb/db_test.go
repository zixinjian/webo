package sdb

import (
	"fmt"
	"os"
	"testing"
)

func TestDbQuery(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)

	db, err := OpenDb(tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	_, _ = db.Db.Exec("CREATE TABLE user (x int, y float)")
	tables, err := GetTableNames(db.Db)
	fmt.Println("tables:", tables)
	DbMultiExe(db.Db, CreateInsertSql("user", []string{"x", "y"}), [][]interface{}{{"1", float64(2)}})
	m, err := DbQuery(db.Db, "SELECT * FROM user")
	fmt.Println(m)
	x := m[0]["x"]
	switch xx := x.(type) {
	default:
		fmt.Println(fmt.Sprintf("%T", xx))
	}
}