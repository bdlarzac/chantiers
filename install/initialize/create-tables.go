/******************************************************************************
    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-09-26 16:04:00+02:00, Thierry Graff : Creation
********************************************************************************/
package initialize

import (
	"bdl.local/bdl/ctxt"
	"fmt"
	"io/ioutil"
	"path"
)

// crée une table ou un type à partir du contenu d'un fichier .sql
func CreateTable(table string) {
	ctx := ctxt.NewContext()
	db := ctx.DB

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		err = tx.Commit()
		if err != nil {
			panic(err)
		}
	}()

	dirSql := getCreateTableDir()

	filename := path.Join(dirSql, table) + ".sql"
	tmp, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	sql := string(tmp)
	fmt.Printf("Crée table ou type %s\n", table)
	if _, err = tx.Exec("drop table if exists " + table + " cascade"); err != nil {
		panic(err)
	}
	if _, err = tx.Exec("drop type if exists " + table + " cascade"); err != nil {
		panic(err)
	}
	if _, err = tx.Exec(sql); err != nil {
		panic(err)
	}
	//    tx.Commit()
}
