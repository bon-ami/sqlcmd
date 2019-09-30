package main

import (
	"database/sql"
	"github.com/bon-ami/eztools"
	"log"
	"strconv"
	"strings"
)

const notifPref = "Done"

func doQuery(db *sql.DB, rows *sql.Rows) (notif string) {
	notif = notifPref
	col, err := rows.Columns()
	if err != nil {
		log.Println(err)
		return
	}
	rawRes := make([][]byte, len(col))
	dest := make([]interface{}, len(col))
	for i, _ := range rawRes {
		dest[i] = &rawRes[i]
	}
	for rows.Next() {
		notif = notif + "\n"
		err = rows.Scan(dest...)
		if err != nil {
			log.Println(err)
		} else {
			for _, raw := range rawRes {
				if raw != nil {
					notif = notif + " [" + string(raw) + "]"
				}
			}
		}
	}
	return notif
}

func doExec(db *sql.DB, res sql.Result) string {
	var (
		num int64
		err error
	)
	notif := notifPref
	num, err = res.LastInsertId()
	if err != nil {
		log.Println(err)
	} else {
		notif = notif + " LastInsertId=" + strconv.FormatInt(num, 10)
	}
	num, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
	} else {
		notif = notif + " RowsAffected=" + strconv.FormatInt(num, 10)
	}
	return notif
}

func main() {
	var (
		err error
		db  *sql.DB
	)
	db, err = eztools.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var (
		rows  *sql.Rows
		res   sql.Result
		query bool
		notif string
	)
	for {
		str := eztools.PromptStr("Enter SQL command:")
		if len(str) < 1 {
			break
		}
		strU := strings.ToUpper(str)
		if strings.HasPrefix(strU, "DESC ") || strings.HasPrefix(strU, "SELECT ") {
			query = true
			rows, err = db.Query(str)
		} else {
			query = false
			res, err = db.Exec(str)
		}
		if err != nil {
			log.Println(err)
		} else {
			if query {
				if rows != nil {
					notif = doQuery(db, rows)
					rows.Close()
					rows = nil
				}
			} else {
				notif = doExec(db, res)
			}
			eztools.ShowStr(notif + "\n")
		}
	}
}
