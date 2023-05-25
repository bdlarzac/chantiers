/*
    Répare les codes ug mal formattés. 
bdlchantiers=> select code from ug where code like '%.%' order by code;
   code   
----------
 XIX.5
 XIX.77
 XIX.78
 XIX.79
 XIX.80
 XIX.81
 XIX.82
 XIX.83
 XIX.84
 XVIII.31
 XVIII.32
 XVIII.33
 XVIII.34

    ATTENTION, la version courante du fix ne gère pas les cas XIX.5 (id 195) et XIX-5 (id 227)

	Intégration : commit 9453126  2023-05-25 12:21

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
	"strings"
)

func Migrate_2023_05_25_fix_code_ug__21(ctx *ctxt.Context) {
	db := ctx.DB
    stmt_select, err := db.Prepare("select id, code from ug where code like '%.%'")
    if err != nil {
        panic(err)
    }
    defer stmt_select.Close()
    //
	stmt_update, err := db.Prepare("update ug set code=$1 where id=$2")
    if err != nil {
        panic(err)
    }
    defer stmt_update.Close()
    //
    rows, err := stmt_select.Query()
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    var id int
    var code string
    for rows.Next() {
        err := rows.Scan(&id, &code)
        if err != nil {
            panic(err)
        }
        tmp := strings.Split(code, ".")
        newCode := strings.Join(tmp, "-")
fmt.Printf("id=%d - code=%s - new = %s\n",id, code, newCode)
		_, err = stmt_update.Exec(newCode, id)
        if err != nil {
            panic(err)
        }
    }
	fmt.Println("Migration effectuée : 2023-05-25-fix-code-ug--21")
}

