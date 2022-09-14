/******************************************************************************

    Ajout de typeunite.NP
    NP pour nombre de piquets

    Intégration : commit b31602d 2022-02-07 23:28:51 +0100

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2022-02-07 22:37:42+01:00, Thierry Graff : Creation
********************************************************************************/
package dbmigrate

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2022_02_07_unite_piquets(ctx *ctxt.Context) {
	db := ctx.DB
	query := `ALTER TYPE typeunite ADD VALUE 'NP' AFTER 'TO'`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2022-02-07-unite-piquets")
}
