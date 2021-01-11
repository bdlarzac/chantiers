/******************************************************************************

    Lit un fichier csv et le charge dans une map.
    Les clés sont tirées de la première ligne.

    @license    GPL
    @history    2019-11-05 05:36:35+01:00, Thierry Graff : Creation
    @history    2020-12-30 17:40:23+01:00, Luka Peschke : Supprime dépendance vers github.com/recursionpharma/go-csv-map
********************************************************************************/

package tiglib

import (
	"encoding/csv"
	"fmt"
	"os"
)

/**
    @param  sep Column separator
**/
func CsvMap(filename string, sep rune) ([]map[string]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(fd)
	reader.Comma = sep
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	recordsLen := len(records)
	if recordsLen < 2 {
		return nil, fmt.Errorf("at least two records are required, got %d", recordsLen)
	}
	head := records[0]
	headLen := len(head)
	if headLen < 1 {
		return nil, fmt.Errorf("expected at least one column, got %d", headLen)
	}
	res := make([]map[string]string, len(records)-1)
	for i, record := range records[1:] {
		if len(record) != headLen {
			return nil, fmt.Errorf("expected %d columns, got %d", headLen, len(record))
		}
		tmp := make(map[string]string)
		for idx, field := range head {
			tmp[field] = record[idx]
		}
		res[i] = tmp
	}
	return res, nil
}
