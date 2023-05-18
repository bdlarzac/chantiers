/*
*******************************************************************************

	    https://github.com/bdlarzac/chantiers/issues/19
	    Adapte les types postgres aux évolutions du logiciel
	    = supprime les types existant et utilise des char à la place
	    
		Intégration : commit

		@copyright  BDL, Bois du Larzac
		@license    GPL
		@history    2023-05-18 06:45:23+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_05_18_clean_types__19(ctx *ctxt.Context) {
//	drop_types_2023_05_18(ctx)
	drop_tables_2023_05_18(ctx)
	fmt.Println("Migration effectuée : 2023-05-18-clean-types--19")
}

func drop_types_2023_05_18(ctx *ctxt.Context) {
	db := ctx.DB
	queries1 := []string{
	    `alter table plaq alter column essence         type char(2)`,
	    `alter table plaq alter column granulo         type char(3)`,
	    `alter table plaq alter column exploitation    type char(1)`,
	    //
	    `alter table plaqop alter column unite         type char(2)`,
	    `alter table plaqop alter column typop         type char(2)`,
	    //
	    `alter table chautre alter column typevalo     type char(2)`,
	    `alter table chautre alter column typevente    type char(3)`,
	    `alter table chautre alter column exploitation type char(1)`,
	    `alter table chautre alter column essence      type char(2)`,
	    `alter table chautre alter column unite        type char(2)`,
	    //
	    `alter table chaufer alter column exploitation type char(1)`,
	    `alter table chaufer alter column essence      type char(2)`,
	    `alter table chaufer alter column unite        type char(2)`,
	    //
	    `alter table stockfrais alter column typefrais type char(2)`,
	}
	queries2 := []string{
        `drop type typegranulo`,
        `drop type typeop`,
        `drop type typessence`,
        `drop type typestockfrais`,
        `drop type typeunite`,
        `drop type typevalorisation`,
        `drop type typevente`,
        `drop type typexploitation`,
	}
	for _, query := range(queries1){
	    fmt.Println(query)
        _, _ = db.Exec(query)
	}
	for _, query := range(queries2){
	    fmt.Println(query)
        _, _ = db.Exec(query)
	}
}

func drop_tables_2023_05_18(ctx *ctxt.Context) {
//	db := ctx.DB
	queries := []string{
        `drop table essence`,
        `drop table role`,
        `drop table typo`,
	}
	for _, query := range(queries){
	    fmt.Println(query)
//        _, _ = db.Exec(query)
	}
}