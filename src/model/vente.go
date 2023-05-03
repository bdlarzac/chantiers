/*
*****************************************************************************

		Vente générique
	    - vente de plaquettes forestières
	    - vente issue d'un chantier autres valorisations

		@copyright  BDL, Bois du Larzac.
		@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
		@history    2023-05-03 08:14:58+02:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

type Vente struct {
	Id           int
	TypeVente    string // "plaq" ou "autre"
	Titre        string
	URL          string // Chaîne vide ou URL du détail de l'entité, ex "/plaq/32"
	DateVente time.Time
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
	//
	Details interface{}
}

// ************************** Nom *******************************

func (v *Vente) String() string {
	return v.Titre
}


// ************************** Get one *******************************

func GetVente(db *sqlx.DB, typeVente string, idVente int) (vente *Vente, err error) {
	vente = &Vente{}
	vente.TypeVente = typeVente
	switch typeVente {
	case "plaq":
		err = vente.computeOneFromVentePlaq(db, idVente)
		break
	case "autre":
		err = vente.computeOneFromChautre(db, idVente)
		break
	}
	if err != nil {
		return vente, werr.Wrapf(err, "Erreur appel vente.computeOneFrom "+typeVente)
	}
	return vente, nil
}

// ************************** compute one *******************************

func (v *Vente) computeOneFromVentePlaq(db *sqlx.DB, idVente int) (err error) {
	vp, err := GetVentePlaq(db, idVente)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetVentePlaq()")
	}
	err = vp.ComputeClient(db) // Obligatoire pour pouvoir utiliser String()
	if err != nil {
		return werr.Wrapf(err, "Erreur appel VentePlaq.ComputeClient()")
	}
	v.Id = vp.Id
	v.TypeVente = "plaq"
	v.Titre = vp.String()
	v.URL = "/chantier/plaquette/" + strconv.Itoa(idVente)
	v.DateVente = vp.DateVente
	v.TypeValo = "PQ"
	err = vp.ComputeQte(db)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel VentePlaq.ComputeQte()")
	}
	v.Volume = vp.Qte
	v.Unite = "MA"
	v.CodeEssence = "PS"
	v.PrixHT =  vp.PUHT * vp.Qte
	v.PUHT = vp.PUHT
	v.TVA = vp.TVA
	v.NumFacture = vp.NumFacture
	v.DateFacture = vp.DateFacture
	v.Notes = vp.Notes
	return nil
}

func (v *Vente) computeOneFromChautre(db *sqlx.DB, idChautre int) (err error) {
	ch, err := GetChautre(db, idChautre)
	if err != nil {
		return werr.Wrapf(err, "Erreur appel GetChautre()")
	}
	v.Id = ch.Id
	v.TypeVente = "autre"
	v.Titre = ch.Titre
//	v.URL = "/chantier/autre/" + strconv.Itoa(idChautre)
	v.DateVente = ch.DateContrat
	v.TypeValo = ch.TypeValo
	v.Volume = ch.VolumeRealise
	v.Unite = ch.Unite
	v.CodeEssence = ch.Essence
	v.PrixHT =  ch.PUHT * ch.VolumeRealise
	v.PUHT = ch.PUHT
	v.TVA = ch.TVA
	v.NumFacture = ch.NumFacture
	v.DateFacture = ch.DateFacture
	v.Notes = ch.Notes
	return nil
}
