/******************************************************************************
    Remplit la base avec des lieux de stockage de test
    Code pas utilisé en fonctionnement normal.

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-05 16:54:53+01:00, Thierry Graff : Creation
********************************************************************************/
package fixture

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"fmt"
)

// *********************************************************
func FillStockage() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	table := "stockage"
	fmt.Println("Remplit " + table + " avec des données de test")
	stockage := &model.Stockage{Nom: "Hangar de test"}
	id, err := model.InsertStockage(db, stockage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Insertion stockage %d ok\n", id)
}
