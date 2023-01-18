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
	"bdl.dbinstall/bdl/install"
	"fmt"
)

func Migrate_2023_01_16_fix_parcelle(ctx *ctxt.Context) {
	versionSCTL := "2023-01-11"
	install.CreateTable(ctx, "parcelle")
	install.AddParcelleCode11(ctx, versionSCTL)
	install.FillParcelle(ctx, versionSCTL)
	install.FillLiensParcelleLieudit(ctx, versionSCTL)
	install.FillLiensParcelleFermier(ctx, versionSCTL)
	install.FillLiensParcelleUG(ctx)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}
