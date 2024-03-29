/*
	    https://github.com/bdlarzac/chantiers/issues/16
	    Ajouter la notion de rôle aux acteurs

		Intégration : commit 019cab9

		@copyright  BDL, Bois du Larzac
		@license    GPL
		@history    2023-02-24 14:43:07+01:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_04_03_role_acteur__16(ctx *ctxt.Context) {
	create_table_role_2023_04_03(ctx)
	create_table_acteur_role_2023_04_03(ctx)
	fill_table_role_2023_04_03(ctx)
	fmt.Println("Migration effectuée : 2023-04-03-role-acteur--16")
}

func fill_table_role_2023_04_03(ctx *ctxt.Context) {
	db := ctx.DB
	var err error
	codeMap := map[string]string{
		// Chantier plaquettes, opérations simples :
		"PLA-AB": "Abatteur", // = bûcheron
		"PLA-DB": "Débardeur",
		"PLA-DC": "Déchiqueteur",
		"PLA-BR": "Broyeur",
		// Chantier plaquettes, transport :
		"PLT-TR": "Transporteur PF",
		"PLT-CO": "Conducteur transport PF",
		"PLT-PO": "Propriétaire outil transport PF",
		// Chantier plaquettes, rangement :
		"PLR-RG": "Rangeur PF",
		"PLR-CO": "Conducteur rangement PF",
		"PLR-PO": "Propriétaire outil rangement PF",
		// Vente plaquettes :
		"VPL-CL": "Client PF",
		// Vente plaquettes, chargement :
		"VPC-CH": "Chargeur PF",
		"VPC-CO": "Conducteur chargement PF",
		"VPC-PO": "Propriétaire outil chargement PF",
		// Vente plaquettes, livraison
		"VPL-LI": "Livreur PF",
		"VPL-CO": "Conducteur livraison PF",
		"VPL-PO": "Propriétaire outil livraison PF",
		// Chantier autres valorisations
		"AVC-PP": "Client pâte à papier",
		"AVC-CH": "Client bois de chauffage",
		"AVC-PL": "Client palettes",
		"AVC-PI": "Client piquets",
		"AVC-BO": "Client bois d'oeuvre",
		// Divers:
		"DIV-MH": "Mesureur d'humidité",
		"DIV-PF": "Propriétaire foncier",
		"DIV-FO": "Fournisseur de plaquettes",
		"FER-BC": "Fermier bois de chauffage",
	}
	//
	stmt, err := db.Prepare("insert into role(code, nom) values($1,$2)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	//
	for code, nom := range codeMap {
		_, err = stmt.Exec(code, nom)
		if err != nil {
			panic(err)
		}
	}
}

func create_table_acteur_role_2023_04_03(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	//
	query = `drop table if exists acteur_role`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table acteur_role (
            id_acteur             int not null references acteur(id),
            code_role             char(6) not null references role(code),
            primary key(id_acteur, code_role)
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	//
	query = `CREATE INDEX acteur_role_acteur_idx ON acteur_role(id_acteur);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	//
	query = `CREATE INDEX acteur_role_role_idx ON acteur_role(code_role);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func create_table_role_2023_04_03(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error

	query = `drop table if exists role`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table role (
            code        char(6) not null primary key,
            nom         varchar(255) not null
        )
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `CREATE INDEX role_idx ON role(code);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}
