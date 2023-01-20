/******************************************************************************

    Ajoute un lien entre les chantiers (plaquettes et autres valorisations)

    Intégration : commit

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-04 08:53:52+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.dbinstall/bdl/install"
	"fmt"
)

func Migrate_2023_01_23_chantier_parcelle(ctx *ctxt.Context) {
	fmt.Println("ok, ici")
	fmt.Println("install.GetDataDir() = " + install.GetDataDir())
	return
	db := ctx.DB
	query := `alter table venteplaq add column facturelivraisonnbkm numeric not null default 0`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2022-09-24-km-livraison")
}
