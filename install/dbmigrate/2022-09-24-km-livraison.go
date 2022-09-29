/******************************************************************************

    Ajout de venteplaq.facturelivraisonnbkm

    Intégration : commit 

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2022-09-24 19:15:20+02:00, Thierry Graff : Creation
********************************************************************************/
package dbmigrate

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2022_09_24_km_livraison(ctx *ctxt.Context) {
fmt.Println("ok, ici")
return
	db := ctx.DB
	query := `ALTER TABLE venteplaq ADD column facturelivraisonnbkm integer`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2022-09-24-km-livraison")
}
