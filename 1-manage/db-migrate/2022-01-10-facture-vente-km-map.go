/******************************************************************************

    Ajout de venteplaq.facturelivraisonunite
    Peut prendre les valeurs 'km', 'map', ''

    Intégration : commit cbab3d5 - 2022-01-10 18:49:11

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2022-01-10 14:02:05+01:00, Thierry Graff : Creation
********************************************************************************/
package dbmigrate

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2022_01_10_facture_vente_km_map(ctx *ctxt.Context) {
	db := ctx.DB
	query := `alter table venteplaq add column facturelivraisonunite text default ''`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2022-01-10-facture-vente-km-map")
}
