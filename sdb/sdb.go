package sdb

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"
	"wb/cs"
	"wb/ii"
	"wb/lg"
)

type Sdb struct {
	Name      string
	Path      string
	FileName  string
	Db        *sql.DB
	logWriter io.Writer
}
type TableInfo struct {
	Name     string
	ColInfos []ColumnInfo
}
type ColumnInfo struct {
	Name      string
	Type      string
	NotNull   bool
	DfltValue interface{}
}

//var dataBaseCache = &_sdbCache{cache: make(map[string]*Sdb), nameMap: make(map[string]string)}

func OpenDb(dbPath string, names ...string) (*Sdb, error) {
	name := defaultDbName
	fmt.Println("names:", names)
	if len(names) > 0 {
		name = names[0]
	}
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		lg.Error("sdb.OpenDb abs Path error:", err.Error())
		return nil, err
	}
	db, err := sql.Open("sqlite3", absPath)
	if err != nil {
		lg.Error("sdb.OpenDb sql.OpenDb error:", err.Error())
		return nil, err
	}
	dir := filepath.Dir(absPath)
	lg.Info("sdb.OpenDb in:", dir, " Name:", name)
	logDir := filepath.Join(dir, "log")
	w, _ := lg.NewRotateFileWriter(logDir, "sdb")
	ndb := &Sdb{Path: absPath, Name: name, Db: db, logWriter: w}
	return ndb, nil
}
func (sdb *Sdb) GetPath() string {
	return path.Join(sdb.Path, sdb.Name)
}
func (sdb *Sdb) query(query string, args ...interface{}) (*sql.Rows, error) {
	tm := time.Now()
	res, err := sdb.Db.Query(query, args...)
	sdb.logExec("query", query, tm, err, args...)
	return res, err
}
func (sdb *Sdb) QueryValue(query string, args ...interface{}) (ret []cs.MObject, err error) {
	ret, err = DbQuery(sdb.Db, query, args...)
	if err != nil {
		fmt.Println("GetColumns error:", err.Error())
		return ret, err
	}
	return ret, nil
}
func (sdb *Sdb) GetColumns(table string) ([]ColumnInfo, error) {
	return GetColumns(sdb.Db, table)
}
func GetTableNames(db *sql.DB) ([]string, error) {
	query := ` select Name from sqlite_master WHERE type = "table"`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	tables := make([]string, 0)
	for rows.Next() {
		var table sql.NullString
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		if table.String != "" && table.String != "sqlite_sequence" {
			tables = append(tables, table.String)
		}
	}
	return tables, nil
}

func GetColumns(db *sql.DB, table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf("pragma table_info('%s')", table)
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("GetColumns error:", err.Error())
		return nil, err
	}
	colInfos := make([]ColumnInfo, 0)
	for rows.Next() {
		var tmp, name, typ, notnull, dflt sql.NullString
		err := rows.Scan(&tmp, &name, &typ, &notnull, &dflt, &tmp)
		if err != nil {
			fmt.Println("GetColumns error:", err.Error())
			return nil, err
		}
		var isNotNull bool
		if notnull.String == "1" {
			isNotNull = true
		} else {
			isNotNull = false
		}
		fmt.Println("isNotNull", notnull)
		colInfo := ColumnInfo{name.String, typ.String, isNotNull, dflt.String}
		colInfos = append(colInfos, colInfo)
	}
	return colInfos, nil
}

func (sdb *Sdb) GetTableInfoDict() ([]TableInfo, error) {
	tableNames, err := GetTableNames(sdb.Db)
	if err != nil {
		lg.Error("GetTables error:", err.Error())
		return nil, err
	}
	tableInfos := make([]TableInfo, len(tableNames))
	for i, tableName := range tableNames {
		colInfos, err := GetColumns(sdb.Db, tableName)
		if err != nil {
			return tableInfos, err
		}
		tableInfo := TableInfo{tableName, colInfos}
		tableInfos[i] = tableInfo
	}
	return tableInfos, nil
}
func (sdb *Sdb) CreateTable(tblInfo TableInfo) error {
	sqlCreate := `
	create table foo (id integer not null primary key, Name text);
	delete from foo;
	`
	_, err := sdb.Db.Exec(sqlCreate)
	if err != nil {
		lg.Error("Create table error sql:", sqlCreate, "error:", err.Error())
		return err
	}
	return nil
}
func (sdb *Sdb) AddSpecialCols() error {
	tblInfos, err := sdb.GetTableInfoDict()
	if err != nil {
		lg.Error("AddSpecialCols error:", err.Error())
		return err
	}
	for _, tblInfo := range tblInfos {
		if err := sdb.addSpecialColsToTable(tblInfo); err != nil {
			return err
		}
	}
	return nil
}
func (sdb *Sdb) addSpecialColsToTable(tblInfo TableInfo) error {
	lg.Info("addSpecialColsToTable: table Name:", tblInfo.Name)
	col2AlterMap := map[string]string{
		ii.SpecialFieldStatus:     "ALTER TABLE `%s` ADD `_st` TEXT DEFAULT ''",
		ii.SpecialFieldDeleter:    "ALTER TABLE `%s` ADD `_dr` TEXT DEFAULT ''",
		ii.SpecialFieldDeleteTime: "ALTER TABLE `%s` ADD `_dt` TEXT DEFAULT ''",
		ii.SpecialFieldUpdater:    "ALTER TABLE `%s` ADD `_ur` TEXT DEFAULT ''",
		ii.SpecialFieldUpdateTime: "ALTER TABLE `%s` ADD `_ut` TEXT DEFAULT ''",
	}
	colMaps := make(map[string]ColumnInfo)
	for _, colInfo := range tblInfo.ColInfos {
		colMaps[colInfo.Name] = colInfo
	}
	for colName, query := range col2AlterMap {
		if _, ok := colMaps[colName]; !ok {
			query := fmt.Sprintf(query, tblInfo.Name)
			if _, err := sdb.Db.Exec(query); err != nil {
				lg.Error("addSpecialColsToTable:", tblInfo, " error:", err.Error())
				return err
			} else {
				lg.Info("add col :", colName, "to table:", tblInfo.Name)
			}
		}
	}
	return nil
}
func (sdb *Sdb) logExec(operaton, query string, t time.Time, err error, args ...interface{}) {
	sub := time.Now().Sub(t) / 1000
	flag := " OK"
	if err != nil {
		flag = "FAIL"
	}
	con := fmt.Sprintf(" -[%s] - [%s / %11s / %9dus] - [%s]", sdb.Name, flag, operaton, sub, query)
	cons := make([]string, 0, len(args))
	for _, arg := range args {
		cons = append(cons, fmt.Sprintf("%v", arg))
	}
	if len(cons) > 0 {
		con += fmt.Sprintf(" - `%s`", strings.Join(cons, "`, `"))
	}
	if err != nil {
		con += " - " + err.Error()
	}
	_ = lg.RunLogger.Output("[Q] ", con, sdb.logWriter)
}
