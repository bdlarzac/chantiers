/*
*****************************************************************************

		Activité générique - Représente toute activité = entité avec une date et souvent un prix.
		Stocké dans les tables = types d'activité concernée
	        chaufer
	        chautre
	        plaq
	        // plaqop
	        // plaqrange
	        // plaqtrans
	        // ventecharge
	        // ventelivre

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
"fmt"
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

// ************************** Get many *******************************
func ComputeActivitesFromFiltres(db *sqlx.DB, filtres map[string][]string) (result []*Activite, err error) {
fmt.Printf("ComputeActivitesFromFiltres() - filtres = %+v\n", filtres)
	result = []*Activite{}
	//
	// Première sélection, par filtre période
	//
	var tmp []*Activite
	// if plaq dans le filtre activite (pas implémenté)
	tmp, err = computePlaqFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computePlaqFromFiltrePeriode()")
	}
	result = append(result, tmp...)
	//
	tmp, err = computeChautreFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computeChautreFromFiltrePeriode()")
	}
	result = append(result, tmp...)
	//
	tmp, err = computeChauferFromFiltrePeriode(db, filtres["periode"])
	if err != nil {
		return result, werr.Wrapf(err, "Erreur appel computeChauferFromFiltrePeriode()")
	}
	result = append(result, tmp...)
	//
	// Filtres suivants
	//
	if len(filtres["essence"]) != 0 {
		result = filtreActivite_essence(db, result, filtres["essence"])
	}
	//
	if len(filtres["valo"]) != 0 {
		result = filtreActivite_valo(db, result, filtres["valo"])
	}
	//
	if len(filtres["fermier"]) != 0 {
		for _, activite := range result {
			activite.ComputeFermiers(db)
		}
		result = filtreActivite_fermier(db, result, filtres["fermier"])
	}
	//
	if len(filtres["ug"]) != 0 {
		for _, activite := range result {
			activite.ComputeUGs(db)
		}
		result = filtreActivite_ug(db, result, filtres["ug"])
	}
	//
	// préparation (faire le plus tard possible pour optimiser)
	//
	if len(filtres["proprio"]) != 0 || len(filtres["parcelle"]) != 0 {
		for _, activite := range result {
			activite.ComputeLiensParcelles(db)
		}
	}
	//
	if len(filtres["parcelle"]) != 0 {
		result = filtreActivite_parcelle(db, result, filtres["parcelle"])
	}
	//
	if len(filtres["proprio"]) != 0 {
		result, err = filtreActivite_proprio(db, result, filtres["proprio"])
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel filtreActivite_proprio()")
		}
	}
//fmt.Printf("result = %+v\n",result)
	// TODO éventuellement, trier par date
	return result, nil
}

// ****************************************************************************************************
// ************************** Auxiliaires de ComputeActivitesFromFiltres() ****************************
// ****************************************************************************************************

// ************************** Selection initiale, par période *******************************

func computePlaqFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Plaq{}
	query = "select * from plaq"
	if len(filtrePeriode) == 2 {
		query += " where datedeb >= $1 and datedeb <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := plaq2Activite(db, chantier)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel plaq2Activite()")
		}
		result = append(result, tmp)
	}
	return result, nil
}

func computeChautreFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Chautre{}
	query = "select * from chautre"
	if len(filtrePeriode) == 2 {
		query += " where datecontrat >= $1 and datecontrat <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := chautre2Activite(db, chantier)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel chautre2Activite()")
		}
		result = append(result, tmp)
	}
	return result, nil
}

func computeChauferFromFiltrePeriode(db *sqlx.DB, filtrePeriode []string) (result []*Activite, err error) {
	var query string
	chantiers := []*Chaufer{}
	query = "select * from chaufer"
	if len(filtrePeriode) == 2 {
		query += " where datechantier >= $1 and datechantier <= $2"
		err = db.Select(&chantiers, query, filtrePeriode[0], filtrePeriode[1])
	} else {
		err = db.Select(&chantiers, query)
	}
	if err != nil {
		return result, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chantier := range chantiers {
		tmp, err := chaufer2Activite(db, chantier)
		if err != nil {
			return result, werr.Wrapf(err, "Erreur appel chaufer2Activite()")
		}
		result = append(result, tmp)
	}
	return result, nil
}

// ************************** Filtres *******************************
// En entrée : liste d'activités
// En sortie : liste d'activités qui satisfont au filtre

func filtreActivite_essence(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			if a.CodeEssence == f {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

func filtreActivite_valo(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			if a.TypeValo == f {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

func filtreActivite_fermier(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, fermier := range a.Fermiers {
				if fermier.Id == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreActivite_ug(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, ug := range a.UGs {
				if ug.Id == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreActivite_parcelle(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				if lienParcelle.IdParcelle == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result
}

func filtreActivite_proprio(db *sqlx.DB, input []*Activite, filtre []string) (result []*Activite, err error) {
	result = []*Activite{}
	for _, a := range input {
		for _, f := range filtre {
			id, _ := strconv.Atoi(f)
			for _, lienParcelle := range a.LiensParcelles {
				parcelle, err := GetParcelle(db, lienParcelle.IdParcelle)
				if err != nil {
					return result, werr.Wrapf(err, "Erreur appel GetParcelle()")
				}
				if parcelle.IdProprietaire == id {
					result = append(result, a)
					break
				}
			}
		}
	}
	return result, nil
}

// ************************** Conversion de struct vers une Activite *******************************

func plaq2Activite(db *sqlx.DB, ch *Plaq) (a *Activite, err error) {
    a = &Activite{}
	a.Id = ch.Id
	a.Titre = ch.Titre
	a.URL = "/chantier/plaquette/" + strconv.Itoa(a.Id)
	a.DateActivite = ch.DateDebut
	a.TypeValo = "PQ"
	err = ch.ComputeVolume(db)
	if err != nil {
		return a, werr.Wrapf(err, "Erreur appel ComputeVolume()")
	}
	a.Volume = ch.Volume
	a.Unite = "MA"
	a.CodeEssence = ch.Essence
	a.Notes = ch.Notes
	return a, nil
}

func chautre2Activite(db *sqlx.DB, ch *Chautre) (a *Activite, err error) {
    a = &Activite{}
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
	return a, nil
}

func chaufer2Activite(db *sqlx.DB, ch *Chaufer) (a *Activite, err error) {
    a = &Activite{}
	a.Id = ch.Id
	a.Titre = ch.Titre
	//	a.URL = "/chaufer/" + strconv.Itoa(idActivite)
	a.DateActivite = ch.DateChantier
	a.TypeValo = "CF"
	a.Volume = ch.Volume
	a.Unite = ch.Unite
	a.CodeEssence = ch.Essence
	a.Notes = ch.Notes
	return a, nil
}
