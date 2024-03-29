/*
Supprime les fermiers non-agricoles de la base.
L'import SCTL a été modifié par Migrate_2023_01_20_fix_non_agricoles
pour ne plus importer les non-agricoles.
Mais la base n'avait pas été nettoyée.

Voir https://github.com/bdlarzac/chantiers/issues/11
Intégration : commit

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func Migrate_2023_05_22_non_agricoles__20(ctx *ctxt.Context) {
	versionSCTL := "2023-02-02"
	//
	// prepare et check
	//
	idsFermier := computeFermiersNonAgricoles_2023_05_22(ctx, versionSCTL)
	//fmt.Printf("\nids fermier - len = %d : %v\n", len(idsFermier), idsFermier)
	//
	idsParcelle_full := computeParcellesFull_2023_05_22(ctx, idsFermier)
	//fmt.Printf("\nidsParcelle_full - len = %d : %v\n", len(idsParcelle_full), idsParcelle_full)
	//
	idsParcelle_shared := computeSharedParcelles_2023_05_22(ctx, idsFermier, idsParcelle_full)
	//fmt.Printf("\nidsParcelle_shared - len = %d : %v\n", len(idsParcelle_shared), idsParcelle_shared)
	//
	idsParcelle_to_delete := computeParcellesToDelete_2023_05_22(ctx, idsParcelle_full, idsParcelle_shared)
	//fmt.Printf("\nidsParcelle_to_delete - len = %d : %v\n", len(idsParcelle_to_delete), idsParcelle_to_delete)
	//
	// Pour vérifier
	//idsChantier := computeChantiers_2023_05_22(ctx, idsParcelle_to_delete)
	//fmt.Printf("\nidsChantier - len = %d : %v\n", len(idsChantier), idsChantier)
	//
	// clean
	//
	strParcellesToDelete := tiglib.JoinInt(idsParcelle_to_delete, ",")
	clean_parcelle_ug_2023_05_22(ctx, strParcellesToDelete)
	clean_parcelle_lieuDit_2023_05_22(ctx, strParcellesToDelete)
	clean_parcelle_fermier_2023_05_22(ctx, strParcellesToDelete)
	clean_fermiers_2023_05_22(ctx, idsFermier)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}

// ********************************** clean **********************************

// deleted 40 rows in fermier
// avant:  107
// après:  67
func clean_fermiers_2023_05_22(ctx *ctxt.Context, idsFermier []int) {
	strIdsFermiers := tiglib.JoinInt(idsFermier, ",")
	query := "delete from fermier where id in(" + strIdsFermiers + ")"
	res, err := ctx.DB.Exec(query)
	if err != nil {
		panic(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("deleted %d rows in fermier\n", n)
}

// deleted 195 rows in parcelle_fermier
// avant:  2464
// après:  2269
func clean_parcelle_fermier_2023_05_22(ctx *ctxt.Context, strParcellesToDelete string) {
	query := "delete from parcelle_fermier where id_parcelle in(" + strParcellesToDelete + ")"
	res, err := ctx.DB.Exec(query)
	if err != nil {
		panic(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("deleted %d rows in parcelle_fermier\n", n)
}

// deleted 184 rows in parcelle_lieudit
// avant:   2351
// après:   2167
func clean_parcelle_lieuDit_2023_05_22(ctx *ctxt.Context, strParcellesToDelete string) {
	query := "delete from parcelle_lieudit where id_parcelle in(" + strParcellesToDelete + ")"
	res, err := ctx.DB.Exec(query)
	if err != nil {
		panic(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("deleted %d rows in parcelle_lieudit\n", n)
}

// deleted 53 rows in parcelle_ug
// avant:   1516
// après:   1463
func clean_parcelle_ug_2023_05_22(ctx *ctxt.Context, strParcellesToDelete string) {
	query := "delete from parcelle_ug where id_parcelle in(" + strParcellesToDelete + ")"
	res, err := ctx.DB.Exec(query)
	if err != nil {
		panic(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Printf("deleted %d rows in parcelle_ug\n", n)
}

// ********************************** check **********************************

// Liste les chantiers reliés à des parcelles supprimées
// Résultat : aucun chantier
func computeChantiers_2023_05_22(ctx *ctxt.Context, idsParcelle []int) []int {
	res := []int{}
	inStr := tiglib.JoinInt(idsParcelle, ",")
	query := "select id_chantier from chantier_parcelle where id_parcelle in(" + inStr + ")"
	err := ctx.DB.Select(&res, query)
	if err != nil {
		fmt.Println("Erreur query : " + query)
		panic(err)
	}
	sort.Ints(res)
	return res
}

// Calcule les ids parcelles présents dans idsParcelle_full mais pas dans idsParcelle_shared
// Résultat : 184 parcelles
func computeParcellesToDelete_2023_05_22(ctx *ctxt.Context, idsParcelle_full, idsParcelle_shared []int) []int {
	res := []int{}
	for _, elt := range idsParcelle_full {
		if !tiglib.InArray(elt, idsParcelle_shared) {
			res = append(res, elt)
		}
	}
	sort.Ints(res)
	return res
}

// Vérifie s'il existe des parcelles reliées aux fermiers non-agricoles
// qui sont aussi reliées avec des fermiers agricoles
// Résultat : 33 parcelles
func computeSharedParcelles_2023_05_22(ctx *ctxt.Context, idsFermier, idsParcelle []int) []int {
	res := []int{}
	strFermiers := tiglib.JoinInt(idsFermier, ",")
	strParcelles := tiglib.JoinInt(idsParcelle, ",")
	query := "select id_parcelle from parcelle_fermier where id_parcelle in(" + strParcelles + ") and id_fermier not in(" + strFermiers + ")"
	err := ctx.DB.Select(&res, query)
	if err != nil {
		fmt.Println("Erreur query : " + query)
		panic(err)
	}
	res = tiglib.ArrayUnique(res)
	sort.Ints(res)
	return res
}

// Résultat : 217 parcelles
func computeParcellesFull_2023_05_22(ctx *ctxt.Context, idsFermier []int) []int {
	res := []int{}
	inStr := tiglib.JoinInt(idsFermier, ",")
	query := "select id_parcelle from parcelle_fermier where id_fermier in(" + inStr + ")"
	err := ctx.DB.Select(&res, query)
	if err != nil {
		fmt.Println("Erreur query : " + query)
		panic(err)
	}
	res = tiglib.ArrayUnique(res)
	sort.Ints(res)
	return res
}

// Résultat : 44 fermiers
func computeFermiersNonAgricoles_2023_05_22(ctx *ctxt.Context, versionSCTL string) []int {
	dirname := install.GetSCTLDataDir(versionSCTL)
	filename := dirname + string(os.PathSeparator) + "Exploita.csv"
	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		panic("Erreur de lecture de Exploita.csv avec tiglib.CsvMap()")
	}
	res := []int{}
	for _, record := range records {
		idExploitant, err := strconv.Atoi(record["IdExploitant"])
		if idExploitant == 1 {
			continue // "PERSONNE"
		}
		if err != nil {
			panic("Erreur appel strconv.Atoi(%s) pour IdExploitant " + record["IdExploitant"])

		}
		if record["Agricole"] == "0" {
			res = append(res, idExploitant)
		}
	}
	sort.Ints(res)
	return res
}
