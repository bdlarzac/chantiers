package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-adodb"
)

func main() {
	fmt.Println("Fill table communes")

	mdbFile := "dbterres.mdb"
	dsn := fmt.Sprintf("Provider=Microsoft.Jet.OLEDB.4.0;Data Source=%s;", mdbFile)
	dbAccess, err := sql.Open("adodb", dsn)
	if err != nil {
		panic(err)
	}
	defer dbAccess.Close()

	var nom string
	rows, err := dbAccess.Query("select NOM from Commune")
	//rows, err := dbAccess.Query("select IdLieuDit,Libelle from LieuDit")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&nom)
		if err != nil {
			panic(err)
		}
		fmt.Println(nom)
	}
}
