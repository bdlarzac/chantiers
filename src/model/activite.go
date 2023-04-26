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
	TypeValo     string
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
		"chaufer": "Chauffage fermier",
		"chautre": "Autre valorisation",
		"plaq":    "Ch. plaquettes",
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

func (a *Activite) ComputeLiensParcelles(db *sqlx.DB) (err error) {
	a.LiensParcelles, err = computeLiensParcellesOfChantier(db, a.TypeActivite, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeLiensParcellesOfChantier()")
	}
	return nil
}

func (a *Activite) ComputeFermiers(db *sqlx.DB) (err error) {
	a.Fermiers, err = computeFermiersOfChantier(db, a.TypeActivite, a.Id)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel computeFermiersOfChantier()")
	}
	return nil
}

func (a *Activite) ComputeUGs(db *sqlx.DB) (err error) {
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
	}
	if err != nil {
		return activ, werr.Wrapf(err, "Erreur appel activ.computeOneFrom "+typeActivite)
	}
	return activ, nil
}

// ************************** Compute one *******************************

func (a *Activite) computeOneFromPlaq(db *sqlx.DB, idActivite int) (err error) {
	ch, err := GetPlaq(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetPlaq()")
	}
	a.Id = ch.Id
	a.Titre = ch.Titre
	a.URL = "/chantier/plaquette/" + strconv.Itoa(idActivite)
	a.DateActivite = ch.DateDebut
	a.TypeValo = "PQ"
	err = ch.ComputeVolume(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel ComputeVolume()")
	}
	a.Volume = ch.Volume
	a.Unite = "MA"
	a.CodeEssence = ch.Essence
	a.Notes = ch.Notes
	return nil
}

func (a *Activite) computeOneFromChautre(db *sqlx.DB, idActivite int) (err error) {
	ch, err := GetChautre(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChautre()")
	}
	a.Id = ch.Id
	a.Titre = ch.Titre
	//	a.URL = "/chautre/" + strconv.Itoa(idActivite)
	a.DateActivite = ch.DateContrat
	a.Volume = ch.VolumeRealise
	a.TypeValo = ch.TypeValo
	a.Unite = ch.Unite
	a.CodeEssence = ch.Essence
	a.PUHT = ch.PUHT
	a.PrixHT = ch.PUHT * ch.VolumeRealise
	a.TVA = ch.TVA
	a.NumFacture = ch.NumFacture
	a.DateFacture = ch.DateFacture
	a.Notes = ch.Notes
	return nil
}

func (a *Activite) computeOneFromChaufer(db *sqlx.DB, idActivite int) (err error) {
	ch, err := GetChaufer(db, idActivite)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChaufer()")
	}
	a.Id = ch.Id
	a.Titre = ch.Titre
	//	a.URL = "/chaufer/" + strconv.Itoa(idActivite)
	a.DateActivite = ch.DateChantier
	a.TypeValo = "CF"
	a.Volume = ch.Volume
	a.Unite = ch.Unite
	a.CodeEssence = ch.Essence
	a.Notes = ch.Notes
	return nil
}
