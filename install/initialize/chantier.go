/******************************************************************************
    Initialisation des chantiers à partir de feuilles Excel
    Doc :
    https://godoc.org/github.com/tealeg/xlsx
    https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc
    
    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-09-24 12:51:45+02:00, Thierry Graff
    
********************************************************************************/

package initialize

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	// "bdl.local/bdl/ctxt"
	// "bdl.local/bdl/generic/tiglib"
	"github.com/tealeg/xlsx/v3"
)

var pDate = regexp.MustCompile(`\d{2}-\d{2}-\d{2}`)
var pDate2 = regexp.MustCompile(`(\d{2}/\d{2}/\d{4})\D+(\d{2}/\d{2}/\d{4})`) // "09/08/2013 - 10/08/2013"

// *********************************************************
func FillChantierPlaquette() {
    
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
        // date DÉCHIQUETTAGE
        cell, err = sh.Cell(2, 7) // H3
        if err != nil {
            panic(err)
        }
        dates := cell.String()
        dateDeb, dateFin := parseDate(dates)
        fmt.Println(dateDeb, dateFin)
break
    }
    
}


// *********************************************************
// Renvoie un tablea" avec 2 éléments : date début et date fin (au format YYYY-MM-DD)
// Quelle merde ces fichiers excel
// si la cell vaut "20/07/2012", on récupère "07-20-12"
// si la cell vaut "09/08/2013 - 10/08/2013", on récupère "09/08/2013 – 10/08/2013"
func parseDate(s string) (string, string) {
fmt.Println("s =", s)
    dateDeb := ""
    dateFin := ""
    dates := pDate.FindAllString(s, -1)
fmt.Printf("dates = %q\n", dates)
    if len(dates) == 1{
        // une seule date - MM/JJ/AA
        tmp := strings.Split(dates[0], "-")
fmt.Println("tmp =", tmp)
        dateDeb = "20" + tmp[2] + "-" + tmp[0] + "-" + tmp[1]
        dateFin = dateDeb
    } else {
//        dates = pDate2.FindStringSubmatch(-1)
fmt.Printf("2 dates = %q\n", dates)
fmt.Println("A IMPLEMENTER")
    }
    return dateDeb, dateFin
}