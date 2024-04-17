/*
Fixe certains codes UG pas reliés à des parcelles
Voir https://github.com/bdlarzac/chantiers/issues/27

Ce fix est effectéué sur les ugs du PSG1

Intégration: commit

Voir : 
- appli/manage/db-install/install/ug
- appli/manage/data/ug.csv

1 - Vérification table ug

Attention à issue #21 (correction de codes ugs invalides) : nb des ugs en base différent du nb juste après la construction ?
NON car sur la base actuelle : select count(*) from ug;                                 => 659
En réexécutant un version modifiée de appli/manage/db-install/install/ug.FillUG()       => len(records2) = 659
==> OK, aucune migration n'a modifié le nb d'ugs depuis le départ
==> on peut compter sur la stabilité des ids ug.

2 - Vérification table parcelle_ug

Sur la base actuelle : select count(*) from parcelle_ug;                                        => 1516

En réexécutant une version modifiée de appli/manage/db-install/install/ug.FillLiensParcelleUG()  
et en comptant le nb de fois que stmt.Exec(id_parcelle, id_ug) est exécuté                      => 1784

DONC PROBLEME : ne pas comparer avec le code de création
mais voir issue #11 appli/manage/db-migrate/2023-01-16-fix-parcelle.go
En réexécutant une version modifiée de fillLiensParcelleUG_2023_01_16()
et en comptant le nb de fois que stmt.Exec(idParcelle, idUG) est exécuté                        => 1528

Mais voir issue #11 
parcelle_ug est aussi modifiée par 2023-05-22-non-agricoles--20.go
Dans commentaire de clean_parcelle_ug_2023_05_22()
// deleted 53 rows in parcelle_ug
// avant:   1516
// après:   1463

==> ???
- pourquoi la base actuelle a 1516 et pas 1463 ?
- pourquoi le commentaire de clean_parcelle_ug_2023_05_22() indique avant = 1516 et pas avant = 1784 ?

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2024-04-17 08:35:38+02:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/generic/tiglib"
	"path"
	"fmt"
)

func Migrate_2024_04_17_ugs_orphelines__27(ctx *ctxt.Context) {
    analyze_2024_04_17ctx(ctx)
	fmt.Println("Migration effectuée : 2024-04-17-ugs-orphelines")
}

func analyze_2024_04_17ctx(*ctxt.Context) {
	filename := path.Join(install.GetDataDir(), "ug.csv")
	records, _ := tiglib.CsvMap(filename, ',')
	for _, record := range records {
		code_ug := record["PG"]
		if code_ug == "" || code_ug == "0" {
			continue
		}
		code_parcelle6 := record["PC"]
		code_parcelle11 := record["ID_PARCELLE_11"]
// fmt.Printf("%s %s\n",code_parcelle6,code_parcelle11)
// break
//		if (code_parcelle6 == "0" || code_parcelle6 == "") && code_parcelle11 != "" {
//		if code_parcelle6 == "0" || code_parcelle6 == "" {
        if code_ug == "X-22" {
// fmt.Printf("%s %s\n",code_parcelle6,code_parcelle11)
			fmt.Printf("%s %s %s\n", code_ug, code_parcelle6, code_parcelle11)
		}
	}
}