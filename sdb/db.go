package sdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"wb/cs"
	"wb/lg"
	"wb/st"
)

const Q = "'"

func DbQuery(db *sql.DB, query string, values ...interface{}) (mos []cs.MObject, err error) {
	rows, err := db.Query(query, values...)
	if err != nil {
		lg.Error("DbQuery error", err.Error())
		return
	}
	defer rows.Close()
	cnt := 0
	var columns []string
	for rows.Next() {
		if cnt == 0 {
			cnt = 1
			columns, err = rows.Columns()
			if err != nil {
				return
			}
		}
		refs := make([]interface{}, len(columns))
		for i := range columns {
			var ref sql.NullString
			refs[i] = &ref
		}
		if err = rows.Scan(refs...); err != nil {
			return
		}
		mo := make(cs.MObject, len(columns))
		for i, k := range columns {
			ref := refs[i]
			value := reflect.Indirect(reflect.ValueOf(ref)).Interface().(sql.NullString)
			if value.Valid {
				mo[k] = value.String
			} else {
				mo[k] = nil
			}
		}
		mos = append(mos, mo)
	}
	return
}
func DbInsertMos(db *sql.DB, table string, cols []string, mos []cs.MObject) string {
	insertSql := CreateInsertSql(table, cols)
	lstValues := make([][]interface{}, len(mos))
	for i, dict := range mos {
		lstValues[i] = make([]interface{}, len(cols))
		for j, col := range cols {
			if v, ok := dict[col]; ok {
				lstValues[i][j] = v
			} else {
				lstValues[i][j] = nil
			}
		}
	}
	return DbMultiExe(db, insertSql, lstValues)
}
func DbDelete(db *sql.DB, deleteSql string, values ...interface{}) string {
	if _, err := db.Exec(deleteSql, values...); err == nil {
		return st.ErrorDbDelete
	}
	return st.Success
}
func DbMultiExe(db *sql.DB, exeSql string, lstValues [][]interface{}) string {
	tx, err := db.Begin()
	if err != nil {
		lg.Error("Failed to begin transaction:", err, ":", exeSql[:32])
		return st.ErrorDbBeginTx
	}
	s, err := tx.Prepare(exeSql)
	if err != nil {
		lg.Error("Failed to Prepare insert:", err, ":", exeSql[:32])
		return st.ErrorDbPrepare
	}
	defer s.Close()
	for _, values := range lstValues {
		_, err = s.Exec(values...)
		if err != nil {
			lg.Error("Failed to exec:", err, ":", exeSql[:32])
			_ = tx.Rollback()
			return st.ErrorDbExeFailed
		}
	}
	err = tx.Commit()
	if err != nil {
		lg.Error("Failed to commit transaction:", err, ":", exeSql[:32])
		return st.ErrorDbCommit
	}
	return st.Success
}
func CreateInsertSql(table string, names []string) string {
	marks := make([]string, len(names))
	for i := range marks {
		marks[i] = "?"
	}
	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(names, sep)
	return fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, table, Q, Q, columns, Q, qmarks)
}
