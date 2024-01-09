/*
Code commun aux différents types de chantier
Pour gérer les liens entre chantiers et autres entités

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-01-26 11:16:37+01:00, Thierry Graff : Creation à partir de ChauferParcelle
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

// Lien entre une parcelle et un chantier (table chantier_parcelle)
// Utilisé par plaq, chautre, chaufer
// Pour chaque parcelle, on doit préciser s'il s'agit d'une parcelle entière ou pas.
// S'il ne s'agit pas d'une parcelle entière, il faut préciser la surface concernée par l'opération.
type ChantierParcelle struct {
	TypeChantier string `db:"type_chantier"` // "plaq", "chautre" ou "chaufer"
	IdChantier   int    `db:"id_chantier"`
	IdParcelle   int    `db:"id_parcelle"`
	Entiere      bool
	Surface      float64
	// Pas stocké en base
	Parcelle *Parcelle
}

// ************************** Liens chantier parcelle *******************************

func computeLiensParcellesOfChantier(db *sqlx.DB, typeChantier string, idChantier int) (result []*ChantierParcelle, err error) {
	query := `select * from chantier_parcelle where type_chantier='` + typeChantier + `' and id_chantier=$1`
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

func deleteLiensChantierParcelle(db *sqlx.DB, typeChantier string, idChantier int) (err error) {
	query := "delete from chantier_parcelle where type_chantier='" + typeChantier + "' and id_chantier=$1"
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func updateLiensChantierParcelle(db *sqlx.DB, typeChantier string, idChantier int, liensParcelles []*ChantierParcelle) (err error) {
	err = deleteLiensChantierParcelle(db, typeChantier, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierParcelle() à partir de updateLiensChantierParcelle()")
	}
	//
	err = insertLiensChantierParcelle(db, typeChantier, idChantier, liensParcelles)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel insertLiensChantierParcelle() à partir de updateLiensChantierParcelle()")
	}
	return nil
}

// ************************** Liens chantier UG *******************************

func computeUGsOfChantier(db *sqlx.DB, typeChantier string, idChantier int) (result []*UG, err error) {
	query := `select * from ug where id in(
	    select id_ug from chantier_ug where type_chantier='` + typeChantier + `' and id_chantier=$1
    )`
	err = db.Select(&result, query, &idChantier)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query : "+query) // result vide (car fonction appelée que si chantier.UGs n'est pas vide)
	}
	return result, nil
}

func computeIdsChantiersFromUG(db *sqlx.DB, typeChantier string, idUG int) (idsUG []int, err error) {
	query := `select id from ` + typeChantier + ` where id in(
	    select id_chantier from chantier_ug where type_chantier='` + typeChantier + `' and id_ug =$1
    )`
	err = db.Select(&idsUG, query, idUG)
	if err != nil {
		return idsUG, werr.Wrapf(err, "Erreur query : "+query)
	}
	return idsUG, nil
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

func deleteLiensChantierUG(db *sqlx.DB, typeChantier string, idChantier int) (err error) {
	query := "delete from chantier_ug where type_chantier='" + typeChantier + "' and id_chantier=$1"
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func updateLiensChantierUG(db *sqlx.DB, typeChantier string, idChantier int, idsUG []int) (err error) {
	err = deleteLiensChantierUG(db, typeChantier, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierUG() à partir de updateLiensChantierUG()")
	}
	//
	err = insertLiensChantierUG(db, typeChantier, idChantier, idsUG)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel insertLiensChantierUG() à partir de updateLiensChantierUG()")
	}
	return nil
}

// ************************** Liens chantier Lieudits *******************************

func computeLieuditsOfChantier(db *sqlx.DB, typeChantier string, idChantier int) (result []*Lieudit, err error) {
	query := `select * from lieudit where id in(
	    select id_lieudit from chantier_lieudit where type_chantier='` + typeChantier + `' and id_chantier=$1
    )`
	err = db.Select(&result, query, &idChantier)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query : "+query)
	}
	return result, nil
}

func insertLiensChantierLieudit(db *sqlx.DB, typeChantier string, idChantier int, idsLieudit []int) (err error) {
	query := `insert into chantier_lieudit(
        type_chantier,
        id_chantier,
        id_lieudit) values($1,$2,$3)`
	for _, idLieudit := range idsLieudit {
		_, err = db.Exec(
			query,
			typeChantier,
			idChantier,
			idLieudit)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func deleteLiensChantierLieudit(db *sqlx.DB, typeChantier string, idChantier int) (err error) {
	query := "delete from chantier_lieudit where type_chantier='" + typeChantier + "' and id_chantier=$1"
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func updateLiensChantierLieudit(db *sqlx.DB, typeChantier string, idChantier int, idsLieudit []int) (err error) {
	err = deleteLiensChantierLieudit(db, typeChantier, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierLieudit() à partir de updateLiensChantierLieudit()")
	}
	//
	err = insertLiensChantierLieudit(db, typeChantier, idChantier, idsLieudit)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel insertLiensChantierLieudit() à partir de updateLiensChantierLieudit()")
	}
	return nil
}

// ************************** Liens chantier Fermiers *******************************

func computeFermiersOfChantier(db *sqlx.DB, typeChantier string, idChantier int) (result []*Fermier, err error) {
	query := `select * from fermier where id in(
	    select id_fermier from chantier_fermier where type_chantier='` + typeChantier + `' and id_chantier=$1
    )`
	err = db.Select(&result, query, &idChantier)
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query : "+query)
	}
	return result, nil
}

func insertLiensChantierFermier(db *sqlx.DB, typeChantier string, idChantier int, idsFermier []int) (err error) {
	query := `insert into chantier_fermier(
        type_chantier,
        id_chantier,
        id_fermier) values($1,$2,$3)`
	for _, idFermier := range idsFermier {
		_, err = db.Exec(
			query,
			typeChantier,
			idChantier,
			idFermier)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}

func deleteLiensChantierFermier(db *sqlx.DB, typeChantier string, idChantier int) (err error) {
	query := "delete from chantier_fermier where type_chantier='" + typeChantier + "' and id_chantier=$1"
	_, err = db.Exec(query, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	return nil
}

func updateLiensChantierFermier(db *sqlx.DB, typeChantier string, idChantier int, idsFermier []int) (err error) {
	err = deleteLiensChantierFermier(db, typeChantier, idChantier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel deleteLiensChantierLieudit() à partir de updateLiensChantierFermier()")
	}
	//
	err = insertLiensChantierFermier(db, typeChantier, idChantier, idsFermier)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel insertLiensChantierFermier() à partir de updateLiensChantierFermier()")
	}
	return nil
}
