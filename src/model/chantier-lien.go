/******************************************************************************
    Code commun aux différents types de chantier
    Pour gérer les liens entre chantiers et autres entités

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2023-01-26 11:16:37+01:00, Thierry Graff : Creation à partir de ChauferParcelle
********************************************************************************/
package model

import(
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
//    "fmt"
)

/** 
    Lien entre une parcelle et un chantier (table chantier_parcelle)
    Utilisé par plaq, chautre, chaufer
    Pour chaque parcelle, on doit préciser s'il s'agit d'une parcelle entière ou pas.
    S'il ne s'agit pas d'une parcelle entière, il faut préciser la surface concernée par la coupe.
**/
type ChantierParcelle struct {
    TypeChantier string `db:"type_chantier"`
	IdChantier  int `db:"id_chantier"`
	IdParcelle int `db:"id_parcelle"`
	Entiere    bool
	Surface    float64
	// Pas stocké en base
	Parcelle *Parcelle
}

// ************************** Liens chantier parcelle *******************************

// pour factoriser chantier.ComputeLiensParcelles()
func computeChantierLiensParcelles(db *sqlx.DB, typeChantier string, idChantier int) (result []*ChantierParcelle, err error) {
	query := `select * from chantier_parcelle where type_chantier='`+typeChantier+`' and id_chantier=$1`
	err = db.Select(&result, query, idChantier)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for i, lien := range result {
		result[i].Parcelle, err = GetParcelle(db, lien.IdParcelle)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel LienParcelle()")
		}
	}
	return result, nil
}

func insertLiensChantierParcelle(db *sqlx.DB, typeChantier string, idChantier int, liensParcelles []*ChantierParcelle) (err error) {
	query := "insert into chantier_parcelle values($1,$2,$3,$4,$5)"
	for _, lien := range liensParcelles {
		_, err = db.Exec(
			query,
			idChantier,
			lien.IdParcelle,
			lien.Entiere,
			lien.Surface,
			typeChantier)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func updateLiensChantierParcelle(db *sqlx.DB, typeChantier string, idChantier int, liensParcelles []*ChantierParcelle) (err error) {
	query := `delete from chantier_parcelle where type_chantier='`+typeChantier+`' and id_chantier=$1`
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "insert into chantier_parcelle values($1,$2,$3,$4,$5)"
	for _, lien := range liensParcelles {
		_, err = db.Exec(
			query,
			lien.IdChantier,
			lien.IdParcelle,
			lien.Entiere,
			lien.Surface,
			typeChantier)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

// ************************** Liens chantier UG *******************************

// pour factoriser chantier.ComputeUGs()
func computeChantierUGs(db *sqlx.DB, typeChantier string, idChantier int) (result []*UG, err error) {
    query := `select * from ug where id in(
	    select id_ug from chantier_ug where type_chantier='`+typeChantier+`' and id_chantier=$1
    )`
	err = db.Select(&result, query, &idChantier)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query : "+query) // result vide (car fonction appelée que si chantier.UGs n'est pas vide)
	}
	return result, nil
}

func insertLiensChantierUG(db *sqlx.DB, typeChantier string, idChantier int, idsUG []int) (err error) {
	query := `insert into chantier_ug(
        type_chantier,
        id_chantier,
        id_ug) values($1,$2,$3)`
	for _, idUG := range idsUG {
		_, err = db.Exec(
			query,
			typeChantier,
			idChantier,
			idUG)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func updateLiensChantierUG(db *sqlx.DB, typeChantier string, idChantier int, idsUG []int) (err error) {
	query := "delete from chantier_ug where type_chantier='"+typeChantier+"' and id_chantier=$1"
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = `insert into chantier_ug(
        type_chantier,
        id_chantier,
        id_ug) values($1,$2,$3)`
	for _, idUG := range idsUG {
		_, err = db.Exec(
			query,
			typeChantier,
			idChantier,
			idUG)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}
