/*
*****************************************************************************

		Activité générique - Représente toute activité = entité avec une date et souvent un prix.
		Stocké dans les tables = types d'activité concernée
	        chaufer
	        chautre
	        plaq
	        plaqop
	        plaqrange
	        plaqtrans
	        ventecharge
	        ventelivre
	        venteplaq

		@copyright  BDL, Bois du Larzac.
		@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
		@history    2023-03-09 14:54:36+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Activite struct {
	Id           int
	TypeActivite string
	Titre        string
	URL          string // Chaîne vide ou URL du détail de l'entité, ex "/plaq/32"
	DateActivite time.Time
	PrixHT       float64
	PUHT         float64
	TVA          float64
	NumFacture   string
	DateFacture  time.Time
	Notes        string
	//
	Details interface{}
}

// ************************** Nom *******************************

func (ch *Activite) String() string {
	return ch.Titre
}

// ************************** Get one *******************************

/*
 */
func GetActivite(db *sqlx.DB, typeActivite string, idActivite int) (activ *Activite, err error) {
	activ = &Activite{}
	activ.TypeActivite = typeActivite
	switch typeActivite {
	case "chaufer":
		err = activ.computeFromChaufer(db, idActivite)
		break
	case "chautre":
		err = activ.computeFromChautre(db, idActivite)
		break
	case "plaq":
		err = activ.computeFromPlaq(db, idActivite)
		break
	case "plaqop":
		err = activ.computeFromPlaqop(db, idActivite)
		break
	case "plaqrange":
		err = activ.computeFromPlaqrange(db, idActivite)
		break
	case "plaqtrans":
		err = activ.computeFromPlaqtrans(db, idActivite)
		break
	case "ventecharge":
		err = activ.computeFromVentecharge(db, idActivite)
		break
	case "ventelivre":
		err = activ.computeFromVentelivre(db, idActivite)
		break
	case "venteplaq":
		err = activ.computeFromVenteplaq(db, idActivite)
		break
	}
	if err != nil {
		return activ, werr.Wrapf(err, "Erreur appel activ.computeFrom "+typeActivite)
	}
	return activ, nil
}

func (activ *Activite) computeFromChaufer(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetChaufer(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChaufer()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
	activ.URL = "/chaufer/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateChantier
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromChautre(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetChautre(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChautre()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
	activ.URL = "/chautre/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateContrat
	activ.PUHT = a.PUHT
	activ.TVA = a.TVA
	activ.NumFacture = a.NumFacture
	activ.DateFacture = a.DateFacture
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromPlaq(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetPlaq(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaq()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
	activ.URL = "/plaq/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateDebut
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromPlaqop(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetPlaqOp(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaqOp()")
	}
	activ.Id = a.Id
	activ.Titre = LabelActivite(a.TypOp)
	activ.URL = "/plaq/" + strconv.Itoa(a.IdChantier)
	activ.DateActivite = a.DateDebut
	activ.PUHT = a.PUHT
	activ.TVA = a.TVA
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromPlaqrange(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetPlaqRange(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaqRange()")
	}
	activ.Id = a.Id
	activ.Titre = "Rangement plaquettes"
	activ.URL = "/plaq/" + strconv.Itoa(a.IdChantier)
	activ.DateActivite = a.DateRange
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromPlaqtrans(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetPlaqTrans(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaqTrans()")
	}
	activ.Id = a.Id
	activ.Titre = "Transport plaquettes"
	activ.URL = "/vente/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateTrans
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromVentecharge(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetVenteCharge(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVenteCharge()")
	}
	activ.Id = a.Id
	//	activ.URL = "/vente/" + TODO
	activ.Titre = "Chargement plaquette"
	activ.URL = "/vente/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateCharge
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromVentelivre(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetVenteLivre(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVenteLivre()")
	}
	activ.Id = a.Id
	activ.Titre = "Livraison plaquettes"
	activ.URL = "/vente/" + strconv.Itoa(a.IdVente)
	activ.DateActivite = a.DateLivre
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeFromVenteplaq(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetVentePlaq(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVentePlaq()")
	}
	activ.Id = a.Id
	activ.Titre = "Vente plaquettes"
	activ.URL = "/vente/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateVente
	activ.PUHT = a.PUHT
	activ.TVA = a.TVA
	activ.NumFacture = a.NumFacture
	activ.DateFacture = a.DateFacture
	activ.Notes = a.Notes
	return nil
}
