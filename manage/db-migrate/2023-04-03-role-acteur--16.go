/*
*
*****************************************************************************

	    https://github.com/bdlarzac/chantiers/issues/16
	    Ajouter la notion de rôle aux acteurs

		Intégration : commit 

		@copyright  BDL, Bois du Larzac
		@license    GPL
		@history    2023-02-24 14:43:07+01:00, Thierry Graff : Creation

*******************************************************************************
*
*/
package main

import (
//	"bdl.dbinstall/bdl/install"
	"bdl.local/bdl/ctxt"
//	"bdl.local/bdl/model"
	"fmt"
)

func Migrate_2023_04_03_role_acteur__16(ctx *ctxt.Context) {
//	create_table_role_2023_04_03(ctx)
//	create_table_role_acteur_2023_04_03(ctx)
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
        "PLT-TR": "Transporteur",
        "PLT-CO": "Conducteur",
        "PLT-PO": "Propriétaire outil",
        // Chantier plaquettes, rangement :
        "PLR-RG": "Rangeur",
        "PLR-CO": "Conducteur",
        "PLR-PO": "Propriétaire outil",
        // Vente plaquettes :
        "VPL-CL": "Client plaquette",
        // Vente plaquettes, chargement :
        "VPC-CH": "Chargeur",
        "VPC-CO": "Conducteur",
        "VPC-PO": "Propriétaire outil",
        // Vente plaquettes, livraison
        "VPL-LI": "Livreur",
        "VPL-CO": "Conducteur",
        "VPL-PO": "Propriétaire outil",
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
	}
	//
	stmt, err := db.Prepare("insert into role(code, nom) values($1,$2)")
	defer stmt.Close()
	if err != nil {
		panic(err)
	}
	//
	for code, nom := range(codeMap) {
		_, err = stmt.Exec(code, nom)
		if err != nil {
			panic(err)
		}
	}
}

func create_table_role_acteur_2023_04_03(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	var err error
	//
	query = `drop table if exists role_acteur`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	query = `
        create table role_acteur (
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
	query = `CREATE INDEX role_acteur_acteur_idx ON role_acteur(id_acteur);`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	//
	query = `CREATE INDEX role_acteur_role_idx ON role_acteur(code_role);`
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

