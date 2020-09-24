/******************************************************************************
    Initialisation des chantiers à partir de feuilles Excel

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-09-24 12:51:45+02:00, Thierry Graff
    
https://godoc.org/github.com/tealeg/xlsx
https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc
********************************************************************************/

package initialize

import (
	"fmt"
	"path"
	"regexp"
	// "bdl.local/bdl/ctxt"
	// "bdl.local/bdl/generic/tiglib"
	"github.com/tealeg/xlsx/v3"
)

var pDate = regexp.MustCompile(`\d{2}/\d{2}/\d{4}`)

// *********************************************************
func FillChantier() {
    
    wb, err := xlsx.OpenFile(path.Join(getPrivateDir(), "Chantiers plaquettes.xlsx"))
    if err != nil {
        panic(err)
    }
    for i, sh := range wb.Sheets {
        fmt.Println("-----", i, sh.Name)
        var cell *xlsx.Cell
        var err error
        // lieudit
        cell, err = sh.Cell(1, 7) // H2
        if err != nil {
            panic(err)
        }
        lieudit, err := cell.FormattedValue()
        fmt.Println("lieu =", lieudit)
        // date broyage
        cell, err = sh.Cell(2, 7) // H3
        if err != nil {
            panic(err)
        }
        dates := cell.String()
        fmt.Println(dates)
break
    }
    
}


// *********************************************************
    // Quelle merde ces fichiers excel
    // si la cell vaut "20/07/2012", on récupère "07-20-12"
    // si la cell vaut "09/08/2013 - 10/08/2013", on récupère "09/08/2013 – 10/08/2013"
func parseDate(s string) [2]string {
    fmt.Println("s =", s)
    dates := pDate.FindAll([]byte(s), -1)
    fmt.Printf("dates = %q\n", dates)
    
    return [2]string{"", ""}
}