/******************************************************************************
    Initialisation UG.
    Code à exécuter à chaque maj du PSG

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-11-14 10:08:26+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"math"
	"path"
	"strconv"
	"strings"
)

func FillUG(ctx *ctxt.Context) {
	table := "ug"
	csvname := "ug.csv"
	fmt.Println("Remplit table " + table + " à partir de " + csvname)
	dirCsv := GetDataDir()
	filename := path.Join(dirCsv, csvname)
	records, err := tiglib.CsvMap(filename, ',')
	if err != nil {
		panic(err)
	}
	//
	// records2 : contient type_coupe ; plusieurs types_coupes pour un code ug - sert à dédoublonner au sein d'un code ug
	//
	records2 := make(map[string]map[string]string)
	// records 3 : previsionnel_coupe et type_peuplement toujours unique pour un code ug
	records3 := make(map[string]map[string]string)
	for _, record := range records {
		code := record["PG"]
		if code == "" || code == "0" {
			continue
		}
		if _, ok := records2[code]; !ok {
			records2[code] = make(map[string]string)
		}
		type_coupe := record["Coupe"] + " " + record["Annee_intervention"] // TODO voir si on garde info "0 0"
		// Arrondit surface au m2 (4 chiffres après la virgule)
		surface, _ := strconv.ParseFloat(record["Surf_SIG"], 64)
		surface = math.Round(surface*10000) / 10000
		surfaceStr := strconv.FormatFloat(surface, 'f', -1, 64)
		records2[code][type_coupe] = type_coupe
		if _, ok := records3[code]; !ok {
			records3[code] = make(map[string]string)
			records3[code]["previsionnel_coupe"] = record["PSG_suivant"]
			records3[code]["type_peuplement"] = record["Essence"]
			records3[code]["surface_sig"] = surfaceStr
		}
	}
	db := ctx.DB
	for code, record2 := range records2 {
		// array_value() pour faire strings.Join()
		var tmp []string
		for _, tmp2 := range record2 {
			tmp = append(tmp, tmp2)
		}
		type_coupe := strings.Join(tmp, ", ")
		sql := fmt.Sprintf("insert into %s(code,type_coupe,previsionnel_coupe,type_peuplement,surface_sig) values('%s','%s','%s','%s',%s)",
			table,
			code,
			records3[code]["previsionnel_coupe"],
			type_coupe,
			records3[code]["type_peuplement"],
			records3[code]["surface_sig"])
		if _, err = db.Exec(sql); err != nil {
			panic(err)
		}
	}
}

// *********************************************************
// @pre La table parcelle existe et est remplie
// @pre la table ug existe et est remplie
func FillLiensParcelleUG(ctx *ctxt.Context) {
	table := "parcelle_ug"
	fmt.Println("Remplit table " + table + " à partir de ug.csv")
	dirCsv := GetDataDir()
	filename := path.Join(dirCsv, "ug.csv")

	records, err := tiglib.CsvMap(filename, ',')
	if err != nil {
		panic(err)
	}

	records2 := make(map[string]map[string]string) // map[code parcelle][code ug] = code ug
	for _, record := range records {
		code_parcelle := record["ID_PARCELLE_11"]
		code_ug := record["PG"]
		if code_parcelle == "0" {
			continue
		}
		if code_ug == "" || code_ug == "0" {
			continue
		}
		//
		if _, ok := records2[code_parcelle]; !ok {
			records2[code_parcelle] = make(map[string]string)
		}
		records2[code_parcelle][code_ug] = code_ug
	}

	// Assoc codes => ids
	db := ctx.DB
	//	defer db.Close() => ??? si pas commenté, dans le cas de install all, génère un panic dans la fonction suivante installStockage()

	codeIdsParcelles := make(map[string][]int) // code parcelle => ids parcelles
	var (
		id   int
		code string
	)
	rows, err := db.Query("select id,code from parcelle")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &code)
		if err != nil {
			panic(err)
		}
		if _, ok := codeIdsParcelles[code]; !ok {
			codeIdsParcelles[code] = []int{}
		}
		codeIdsParcelles[code] = append(codeIdsParcelles[code], id)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	codeIdUgs := make(map[string]int) // code ug => id ug
	rows, err = db.Query("select id,code from ug")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &code)
		if err != nil {
			panic(err)
		}
		codeIdUgs[code] = id
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	sql := fmt.Sprintf("insert into %s(id_parcelle,id_ug) values($1, $2)", table)
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	for code_parcelle, codes_ug := range records2 {
		for _, code_ug := range codes_ug {
			id_ug := codeIdUgs[code_ug]
			for _, id_parcelle := range codeIdsParcelles[code_parcelle] {
				_, err = stmt.Exec(id_parcelle, id_ug)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
