/******************************************************************************

    Change les codes parcelle : passe de code à 6 caractères (ex 0C0001) à code à 11 caractères (ex 120820C0001)

    Voir https://github.com/bdlarzac/chantiers/issues/11
    Intégration : commit

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"github.com/jmoiron/sqlx"
	"fmt"
	
)

func Migrate_2023_01_16_fix_parcelle(ctx *ctxt.Context) {
	db := ctx.DB
//    migrate_2023_01_16_alter_table_parcelle(db)
//    migrate_2023_01_16_create_table_chantier_parcelle(db)
    migrate_2023_01_16_update(db)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}

// Fonctions auxiliaires (commencent par des minuscules)

func migrate_2023_01_16_alter_table_parcelle(db *sqlx.DB) {
	query := `alter table parcelle alter column code type char(11)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func migrate_2023_01_16_create_table_chantier_parcelle(db *sqlx.DB) {
    var err error
	query := `
	    create table chantier_parcelle(
        type_chantier           varchar(7), -- "plaq" ou "chautre" ou "chaufer"
        id_chantier             int not null,
        id_parcelle             int not null references parcelle(id),
        primary key(type_chantier, id_chantier, id_parcelle)
    )`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX chantier_parcelle_chantier_idx ON chantier_parcelle (type_chantier, id_chantier)`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX chantier_parcelle_parcelle_idx ON chantier_parcelle (type_chantier, id_parcelle)`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

/**
    - Vide la table parcelle et la remplit avec les nouveaux codes
    - Vide la table parcelle_ug et la remplit à nouveau
**/
func migrate_2023_01_16_update(db *sqlx.DB) {
    
}

/* 
select * from chaufer_parcelle
 id_chaufer | id_parcelle | entiere | surface 
------------+-------------+---------+---------
          1 |        1577 | t       |       0
          2 |        1025 | f       |       2
          3 |        1262 | f       |       2
          4 |        1262 | f       |       2
          5 |        1319 | f       |       1
          6 |        1253 | f       |       1
          7 |        2203 | t       |       0

*/

