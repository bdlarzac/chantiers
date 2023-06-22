/*
    Ajoute le certains rôles ne pouvant pas être modifiés par l'interface
    - DIV-PF (propriétaire foncier) aux acteurs SCTL et GFA
    - DIV-FO (fournisseur de plaquettes) à l'acteur BDL
    Ces rôles ne peuvent pas être modifiés via l'interface car la notion de rôle n'existait pas au début
    (code à modifier dans le futur)
    
	@copyright  BDL, Bois du Larzac
	@license    GPL

*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2023_06_21_ajout_roles(ctx *ctxt.Context) {
    db := ctx.DB
    queries := []string{
        "insert into acteur_role(id_acteur,code_role) values(1, 'DIV-PF')", // SCTL
        "insert into acteur_role(id_acteur,code_role) values(2, 'DIV-FO')", // BDL
        "insert into acteur_role(id_acteur,code_role) values(3, 'DIV-PF')", // GFA
    }
    for _, query := range(queries){
        stmt, err := db.Prepare(query)
        defer stmt.Close()
        _, err = stmt.Exec()
        if err != nil {
            panic(err)
        }
    }
	fmt.Println("Migration effectuée : 2023-06-21-ajout-roles")
}
