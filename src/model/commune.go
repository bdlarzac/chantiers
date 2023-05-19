/*
*****************************************************************************

	Communes

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2019-11-07, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
)

type Commune struct {
	Id        int // IdCommune de la base SCTL
	Nom       string
	NomCourt  string
	CodeInsee string
	Lieudits  []*Lieudit
}

const N_COMMUNES = 13

// ************************** Nom *******************************

func (c *Commune) String() string {
	return c.Nom
}

// ************************** Get one *******************************
/*
   Renvoie une Commune contenant Id et Nom.
   Les autres champs ne sont pas remplis.
*/
func GetCommune(db *sqlx.DB, id int) (commune *Commune, err error) {
	commune = &Commune{}
	query := "select * from commune where id=$1"
	row := db.QueryRowx(query, id)
	err = row.StructScan(commune)
	if err != nil {
		return commune, werr.Wrapf(err, "Incapable de fabriquer la commune")
	}
	return commune, nil
}

// ************************** Get many *******************************

/*
Renvoie une liste de communes triés en utilisant un champ de la table
@param field    Champ de la table commune utilisé pour le tri
*/
func GetSortedCommunes(db *sqlx.DB, field string) (communes []*Commune, err error) {
	communes = []*Commune{}
	query := "select * from commune order by " + field
	err = db.Select(&communes, query)
	if err != nil {
		return communes, werr.Wrapf(err, "Erreur query : "+query)
	}
	return communes, nil
}

/*
Renvoie la liste de toutes les communes avec leurs lieux-dits
*/
func ListCommunesEtLieudits(db *sqlx.DB) (communes []*Commune, err error) {
	communes = make([]*Commune, N_COMMUNES)
	query := "select * from commune"
	rows, err := db.Queryx(query)
	if err != nil {
		return communes, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	defer rows.Close()
	queryLD := "select * from lieudit where lieudit.id in(select id_lieudit from commune_lieudit where id_commune=$1)"
	i := 0
	for rows.Next() {
		commune := &Commune{}
		err = rows.StructScan(commune)
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur fabrication Commune")
		}
		rows2, err := db.Queryx(queryLD, commune.Id)
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur query DB : "+query)
		}
		for rows2.Next() {
			lieudit := &Lieudit{}
			err = rows2.StructScan(lieudit)
			if err != nil {
				rows2.Close()
				return communes, werr.Wrapf(err, "Erreur fabrication LieuDit")
			}
			commune.Lieudits = append(commune.Lieudits, lieudit)
		}
		// voir si ok de Close comme ça dans la boucle
		rows2.Close()
		err = rows2.Err()
		if err != nil {
			return communes, werr.Wrapf(err, "Erreur rows2.Next()")
		}
		communes[i] = commune
		i++
	}
	err = rows.Err()
	if err != nil {
		return communes, werr.Wrapf(err, "Erreur rows.Next()")
	}
	return communes, nil
}
