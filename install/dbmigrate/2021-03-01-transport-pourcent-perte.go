/******************************************************************************

    Ajout de plaqtrans.pourcentperte

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2021-03-01 10:08:18+01:00, Thierry Graff : Creation
********************************************************************************/
package dbmigrate

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2021_03_01() {
	ctx := ctxt.NewContext()
	db := ctx.DB
	query := `alter table plaqtrans add column pourcentperte numeric not null`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Migration effectuée : 2021-03-01-transport-pourcent-perte")
}
