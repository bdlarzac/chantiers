/******************************************************************************

    Lit un fichier csv et le charge dans une map.
    Les clés sont tirées de la première ligne.

    @license    GPL
    @history    2019-11-05 05:36:35+01:00, Thierry Graff : Creation
********************************************************************************/

package tiglib

import (
	"github.com/recursionpharma/go-csv-map"
	"os"
)

/**
    @param  sep Column separator
**/
func CsvMap(filename string, sep rune) ([]map[string]string, error) {
	csv, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csvmap.NewReader(csv)
	reader.Reader.Comma = sep
	reader.Columns, err = reader.ReadHeader()
	if err != nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, err
}
