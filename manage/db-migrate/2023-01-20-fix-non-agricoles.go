/*
Supprime les fermiers non-agricoles, importés par erreur lors des maj SCTL précédentes
Note: marqué le 2023-01-20 pour apparaître avant la migration 2023-01-23-chantier-parcelle
Mais en fait écrite le 2023-02-04 alors que la migration 2023-01-23-chantier-parcelle était en cours de développement.
Mais la présente migration a bien été appliquée avant sur la base de prod

Intégration : commit b0ee806

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2023-02-04 08:41:58+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_01_20_fix_non_agricoles(ctx *ctxt.Context) {
	db := ctx.DB
	query := `delete from fermier where id in(1,28,37,42,43,46,48,52,53,54,58,59,63,64,66,69,70,128,129,130,73,74,77,79,85,86,90,91,92,96,97,100,102,103,105,110,113,115,116,117,119,120,121,122,134);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2023-01-20-fix-non-agricole")
}
