/*
*
*****************************************************************************

	Supprime les fermiers non-agricoles de la base.
	L'import SCTL a été modifié par Migrate_2023_01_20_fix_non_agricoles
	pour ne plus importer les non-agricoles.
	Mais la base n'avait pas été nettoyée.

	Voir https://github.com/bdlarzac/chantiers/issues/11
	Intégration : commit 

	@copyright  BDL, Bois du Larzac
	@license    GPL
	@history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation

*******************************************************************************
*
*/
package main

import (
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"fmt"
	"os"
	"strconv"
)

func Migrate_2023_05_22_non_agricoles__20(ctx *ctxt.Context) {
	versionSCTL := "2023-02-02"
	//
	// check
	//
	idsFermier := computeFermiersNonAgricoles_2023_05_22(ctx, versionSCTL)
	//fmt.Printf("ids fermier = \n%v\n",tiglib.JoinInt(idsFermier, ","))
	//
	idsParcelle := computeParcelles_2023_05_22(ctx, idsFermier)
	//fmt.Printf("ids parcelle = \n%v\n",tiglib.JoinInt(idsParcelle, ","))
	//
	idsParcelle_problem := checkParcelles_2023_05_22(ctx, idsFermier, idsParcelle)
	fmt.Printf("idsParcelle_problem = %v\n", idsParcelle_problem)
	//
//	idsChantier := computeChantiers_2023_05_22(ctx, idsParcelle)
	//fmt.Printf("ids chantier = \n%v\n",tiglib.JoinInt(idsChantier, ","))
	//
	// clean
	//
//	strIn := tiglib.JoinInt(idsParcelle, ",")
//	clean_parcelle_ug_2023_05_22(ctx, strIn)
	// clean_parcelle_lieuDit_2023_05_22(ctx, strIn)
	// clean_parcelle_fermier_2023_05_22(ctx, strIn)
	// clean_chantier_parcelle_2023_05_22(ctx, strIn)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}

// ********************************** clean **********************************

func clean_parcelle_ug_2023_05_22(ctx *ctxt.Context, strIn string) {
	query := "delete from parcelle_ug where id_parcelle in(" + strIn + ")"
fmt.Println(query)
	//res, err := db.Exec(query, idChantier)
    //rowCnt, err := res.RowsAffected()
}

// ********************************** check **********************************

// Résultat : aucun chantier
func computeChantiers_2023_05_22(ctx *ctxt.Context, idsParcelle []int) []int {
	idsChantier := []int{}
	inStr := tiglib.JoinInt(idsParcelle, ",")
	query := "select id_chantier from chantier_parcelle where id_parcelle in(" + inStr + ")"
	err := ctx.DB.Select(&idsChantier, query)
	if err != nil {
		fmt.Println("Erreur query : "+query)
	    panic(err)
	}
    return idsChantier
}

// Vérifie si des parcelles reliées aux fermiers non-agricoles
// ne sont pas aussi reliées avec des fermiers agricoles
// Résultat : 37 parcelles
func checkParcelles_2023_05_22(ctx *ctxt.Context, idsFermier, idsParcelle []int) []int {
	res := []int{}
	strFermiers := tiglib.JoinInt(idsFermier, ",")
	strParcelles := tiglib.JoinInt(idsParcelle, ",")
	query := "select id_parcelle from parcelle_fermier where id_parcelle in(" + strParcelles + ") and id_fermier not in(" + strFermiers + ")"
	err := ctx.DB.Select(&res, query)
	if err != nil {
		fmt.Println("Erreur query : "+query)
	    panic(err)
	}
	return res
}

// Résultat : 238 parcelles
func computeParcelles_2023_05_22(ctx *ctxt.Context, idsFermier []int) []int {
	idsParcelle := []int{}
	inStr := tiglib.JoinInt(idsFermier, ",")
	query := "select id_parcelle from parcelle_fermier where id_fermier in(" + inStr + ")"
	err := ctx.DB.Select(&idsParcelle, query)
	if err != nil {
		fmt.Println("Erreur query : "+query)
	    panic(err)
	}
    return idsParcelle
}

// Résultat : 40 fermiers
func computeFermiersNonAgricoles_2023_05_22(ctx *ctxt.Context, versionSCTL string) []int{
    dirname := install.GetSCTLDataDir(versionSCTL)
	filename := dirname + string(os.PathSeparator) + "Exploita.csv"
	records, err := tiglib.CsvMap(filename, ';')
	if err != nil {
		panic("Erreur de lecture de Exploita.csv avec tiglib.CsvMap()")
	}
	result := []int{}
	for _, record := range records {
		idExploitant, err := strconv.Atoi(record["IdExploitant"])
		if idExploitant == 1 {
			continue // "PERSONNE"
		}
		if err != nil {
			panic("Erreur appel strconv.Atoi(%s) pour IdExploitant " + record["IdExploitant"])

		}
		if record["Agricole"] == "0" {
		    result = append(result, idExploitant)
		}
    }
    return result
}
