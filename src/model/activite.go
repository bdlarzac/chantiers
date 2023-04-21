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
	Valorisation string
	Volume       float64
	Unite        string // pour le volume
	CodeEssence  string
	PrixHT       float64
	PUHT         float64
	TVA          float64
	NumFacture   string
	DateFacture  time.Time
	Notes        string
	// relations n-n
	LiensParcelles []*ChantierParcelle
	UGs            []*UG
	Fermiers       []*Fermier
	//
	Details interface{}
}


/*
    Renvoie une map code activité => nom
*/
func GetActivitesMap() map[string]string {
    return map[string]string{
        "chaufer":      "Chauffage fermier",
        "chautre":      "Autre valorisation",
        "plaq":         "Ch. plaquettes",
        // "plaqop":       "",
        // "plaqrange":    "",
        // "plaqtrans":    "",
        // "ventecharge":  "",
        // "ventelivre":   "",
        // "venteplaq":    "",
    }
}

// ************************** Nom *******************************

func (a *Activite) String() string {
	return a.Titre
}

// ************************** Instance methods *******************************

func (a *Activite) ComputeLiensParcelles(db *sqlx.DB) (err error){
    a.LiensParcelles, err = computeLiensParcellesOfChantier(db, a.TypeActivite, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier()")
	}
    return nil
}

func (a *Activite) ComputeFermiers(db *sqlx.DB) (err error){
    a.Fermiers, err = computeFermiersOfChantier(db, a.TypeActivite, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeFermiersOfChantier()")
	}
    return nil
}

func (a *Activite) ComputeUGs(db *sqlx.DB) (err error){
    a.UGs, err = computeUGsOfChantier(db, a.TypeActivite, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeUGsOfChantier()")
	}
    return nil
}

// ************************** Get one *******************************

func GetActivite(db *sqlx.DB, typeActivite string, idActivite int) (activ *Activite, err error) {
	activ = &Activite{}
	activ.TypeActivite = typeActivite
	switch typeActivite {
	case "chaufer":
		err = activ.computeOneFromChaufer(db, idActivite)
		break
	case "chautre":
		err = activ.computeOneFromChautre(db, idActivite)
		break
	case "plaq":
		err = activ.computeOneFromPlaq(db, idActivite)
		break
	/* 
	case "plaqop":
		err = activ.computeOneFromPlaqop(db, idActivite)
		break
	case "plaqrange":
		err = activ.computeOneFromPlaqrange(db, idActivite)
		break
	case "plaqtrans":
		err = activ.computeOneFromPlaqtrans(db, idActivite)
		break
	case "ventecharge":
		err = activ.computeOneFromVentecharge(db, idActivite)
		break
	case "ventelivre":
		err = activ.computeOneFromVentelivre(db, idActivite)
		break
	case "venteplaq":
		err = activ.computeOneFromVenteplaq(db, idActivite)
		break
	*/
	}
	if err != nil {
		return activ, werr.Wrapf(err, "Erreur appel activ.computeOneFrom "+typeActivite)
	}
	return activ, nil
}

// ************************** Compute one *******************************


func (activ *Activite) computeOneFromPlaq(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetPlaq(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaq()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
	activ.URL = "/chantier/plaquette/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateDebut
	activ.Valorisation = "PQ"
	err = a.ComputeVolume(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeVolume()")
	}
	activ.Volume = a.Volume
	activ.Unite = "MA"
	activ.CodeEssence = a.Essence
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeOneFromChautre(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetChautre(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChautre()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
//	activ.URL = "/chautre/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateContrat
	activ.Volume = a.VolumeRealise
	activ.Valorisation = a.TypeValo
	activ.Unite = a.Unite
	activ.CodeEssence = a.Essence
	activ.PUHT = a.PUHT
	activ.TVA = a.TVA
	activ.NumFacture = a.NumFacture
	activ.DateFacture = a.DateFacture
	activ.Notes = a.Notes
	return nil
}

func (activ *Activite) computeOneFromChaufer(db *sqlx.DB, idActivite int) (err error) {
	a, err := GetChaufer(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChaufer()")
	}
	activ.Id = a.Id
	activ.Titre = a.Titre
//	activ.URL = "/chaufer/" + strconv.Itoa(idActivite)
	activ.DateActivite = a.DateChantier
	activ.Valorisation = "CF"
	activ.Volume = a.Volume
	activ.Unite = a.Unite
	activ.CodeEssence = a.Essence
	activ.Notes = a.Notes
	return nil
}

/* 
func (activ *Activite) computeOneFromPlaqop(db *sqlx.DB, idActivite int) (err error) {
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

// supprimer les fonctions suivantes ?

func (activ *Activite) computeOneFromPlaqrange(db *sqlx.DB, idActivite int) (err error) {
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

func (activ *Activite) computeOneFromPlaqtrans(db *sqlx.DB, idActivite int) (err error) {
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

func (activ *Activite) computeOneFromVentecharge(db *sqlx.DB, idActivite int) (err error) {
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

func (activ *Activite) computeOneFromVentelivre(db *sqlx.DB, idActivite int) (err error) {
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

func (activ *Activite) computeOneFromVenteplaq(db *sqlx.DB, idActivite int) (err error) {
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
*/
