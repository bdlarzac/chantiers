/*
@copyright  BDL, Bois du Larzac
@license    GPL
@history    2019-09-26 16:04:00+02:00, Thierry Graff : Creation
*/
package install

import (
	"bdl.local/bdl/ctxt"
	"fmt"
	"io/ioutil"
	"path"
)

// crée une table ou un type à partir du contenu d'un fichier .sql
func CreateTable(ctx *ctxt.Context, table string) {
	db := ctx.DB
	var err error

	dirSql := GetCreateTableDir()

	filename := path.Join(dirSql, table) + ".sql"
	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	sql := string(tmp) // contient l'instruction create table ou create type
	fmt.Printf("Crée table ou type %s\n", table)
	if _, err = db.Exec("drop table if exists " + table + " cascade"); err != nil {
		panic(err)
	}
	if _, err = db.Exec("drop type if exists " + table + " cascade"); err != nil {
		panic(err)
	}
	if _, err = db.Exec(sql); err != nil {
		panic(err)
	}
}
