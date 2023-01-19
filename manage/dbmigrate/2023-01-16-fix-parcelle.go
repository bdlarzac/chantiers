/******************************************************************************

    Change les codes parcelle : passe de code à 6 caractères (ex 0C0001) à code à 11 caractères (ex 120820C0001)

    Voir https://github.com/bdlarzac/chantiers/issues/11
    Intégration : commit

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2023-01-16 15:47:07+01:00, Thierry Graff : Creation
********************************************************************************/
package main

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/tiglib"
	"bdl.dbinstall/bdl/install"
	"path"
	"strconv"
	"fmt"
)

func Migrate_2023_01_16_fix_parcelle(ctx *ctxt.Context) {
	versionSCTL := "2023-01-11"
//	alterTableParcelle(ctx)
//	fillTableParcelle(ctx, versionSCTL)
//	 alterTableCommune(ctx)
//	fillTableCommune(ctx, versionSCTL)
	fillLiensParcelleUG(ctx, versionSCTL)
	fmt.Println("Migration effectuée : 2023-01-16-fix-parcelle")
}

// Fonctions auxiliaires, commencent par des minuscules

func alterTableParcelle(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table parcelle add column id_commune int`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table parcelle add constraint parcelle_id_commune_fkey foreign key(id_commune) references commune(id) not valid`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fillTableParcelle(ctx *ctxt.Context, versionSCTL string) {
	var dirCsv, filename, csvname string
	var err error
	// Charge parcelles
	csvname = "Parcelle.csv"
	dirCsv = install.GetSCTLDataDir(versionSCTL)
	filename = path.Join(dirCsv, csvname)
	parcelles, _ := tiglib.CsvMap(filename, ';')
	// update table parcelle
	db := ctx.DB
	if err != nil {
		panic(err)
	}
	stmt, err := db.Prepare("update parcelle set id_commune=$1 where id=$2")
	for _, parcelle := range parcelles {
        idCommune, _ := strconv.Atoi(parcelle["IdCommune"])
        idParcelle, _ := strconv.Atoi(parcelle["IdParcelle"])
		_, err = stmt.Exec(idCommune, idParcelle)
		if err != nil {
			panic(err)
		}
	}	
}

func alterTableCommune(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	query = `alter table commune add column codeinsee char(5) not null default ''`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `alter table commune alter column codeinsee drop default`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func fillTableCommune(ctx *ctxt.Context, versionSCTL string) {
 	var dirCsv, filename, csvname string
	csvname = "commune.csv"
	dirCsv = install.GetDataDir()
	filename = path.Join(dirCsv, csvname)
	communes, _ := tiglib.CsvMap(filename, ';')
    db := ctx.DB
	stmt, err := db.Prepare("update commune set codeinsee=$1 where id=$2")
	for _, commune := range communes {
		_, err = stmt.Exec(commune["code_insee"], commune["id"])
        if err != nil {
            panic(err)
        }
	}
}

func fillLiensParcelleUG(ctx *ctxt.Context, versionSCTL string) {
	var dirCsv, filename, csvname string
	var err error
	db := ctx.DB
    var query string
	
	// Prépare code parcelle 11 => id parcelle
	parcelle11_id := make(map[string]string, 2351) // 2351 = nb de parcelles en base
	query = "select concat(c.codeinsee, p.code) as code11, p.id as idParcelle from commune as c,parcelle as p where c.id=p.id_commune"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var idParcelle int
	var code11 string
	for rows.Next() {
        err := rows.Scan(&code11, &idParcelle)
        if err != nil {
            panic(err)
        }
        parcelle11_id[code11] = strconv.Itoa(idParcelle)
	}
	
	// Prépare code ug => id ug
	ug_code_id := make(map[string]string, 659) // 659 = nb de ugs en base
	query = "select id,code from ug limit 10"
	rows, err = db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var idUG int
	var codeUG string
	for rows.Next() {
        err := rows.Scan(&idUG, &codeUG)
        if err != nil {
            panic(err)
        }
        ug_code_id[codeUG] = strconv.Itoa(idUG)
	}
	
	// Charge UGs
	csvname = "ug.csv"
	dirCsv = install.GetDataDir()
	filename = path.Join(dirCsv, csvname)
	ugs, err := tiglib.CsvMap(filename, ',') // []map[string]string
	if err != nil {
		panic(err)
	}
	
	// Vide table parcelle_ug
	query = `truncate table parcelle_ug`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	
	// Remplit parcelle_ug en utilisant code à 11 caractères
	for _, ug := range(ugs) {
fmt.Printf("%+v\n",ug)
break
	}
    
}


/** 
    Modifie Parcelle.csv an ajoutant une colonne code11 = code initial (6 caractères) précédé du code INSEE de la commune
**/
/* 
func AddParcelleCode11(ctx *ctxt.Context, versionSCTL string) {
	fmt.Println("Modifie Parcelle.csv (ajoute colonne PARCELLE11)")
	var dirCsv, filename, csvname string
	var records []map[string]string
	
	// Prépare map id commune => code INSEE
	csvname = "commune.csv"
	dirCsv = GetDataDir()
	filename = path.Join(dirCsv, csvname)
	records, _ = tiglib.CsvMap(filename, ';')
	id_insee := make(map[string]string)
	for _, v := range records {
	    id_insee[v["id"]] = v["code_insee"]
	}
	
	// Charge le fichier à modifier
	csvname = "Parcelle.csv"
	dirCsv = GetSCTLDataDir(versionSCTL)
	filename = path.Join(dirCsv, csvname)
	records, _ = tiglib.CsvMap(filename, ';')
	
	// Génère le fichier modifié
	keys := []string{"PARCELLE", "SURFACE", "REVENU", "SCTL", "IdParcelle", "IdGfa", "IdCommune", "IdLieuDit", "IdTypeCad", "IdClassCad", "PARCELLE11"}
	res := strings.Join(keys, ";") + "\n"
	for _, v := range records {
	    for _, key := range(keys[:len(keys)-1]){
            res += v[key] + ";"
	    }
	    res += id_insee[v["IdCommune"]] + v["PARCELLE"] + "\n"
	}
	file, err := os.Create(filename)
	if err != nil {
	    panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(res)
	if err != nil {
	    panic(err)
	}
}

*/