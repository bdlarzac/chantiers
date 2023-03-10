/*
*****************************************************************************

	Calculs utilisés dans les bilans

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2021-01-19 10:09:42+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

/*
    Structure de données adaptée au bilan par valorisation et par essence
    Normalement, aurait dû être :
    type Valorisation map[string]map[string][2]float64
        "BO":
            "CH": {<volume>, <chiffe affaire>}
    
    mais pas été foutu de faire fonctionner ça, donc fait une map du style :
    "BO-CH-vol": <volume>,
    "BO-CH-ca": <chiffe affaire>
*/
type Valorisations map[string]float64

/*
    Renvoie un tableau contenant les dates de début / fin des "saisons"
    Les saisons encadrent tous les chantiers stockés en base.
    Une saison dure un an.
    
    @param limiteSaison string au format JJ/MM (tiré de 'debut-saison' en conf)
    
    @return
      - un tableau de 2 time.Time avec les dates limites des saisons
      - un bool indiquant s'il existe des chantiers en base
      - une erreur éventuelle
*/
func ComputeLimitesSaisons(db *sqlx.DB, limiteSaison string) ([][2]time.Time, bool, error) {
	// retour
	var res [][2]time.Time
	var err error
	//
	// first, last = date du premier et dernier chantier plaquettes en base
	//
	var first, last time.Time
	var query string
	// chantiers plaquettes
	var first1, last1 time.Time
	ok1 := true
	query = "select min(datedeb), max(datedeb) from plaq"
	err = db.QueryRow(query).Scan(&first1, &last1)
	if err != nil {
		ok1 = false
	}
	// chantiers autres valorisations
	ok2 := true
	var first2, last2 time.Time
	query = "select min(datecontrat), max(datecontrat) from chautre"
	err = db.QueryRow(query).Scan(&first2, &last2)
	if err != nil {
		ok2 = false
	}
	// chantiers chauffage fermier
	ok3 := true
	var first3, last3 time.Time
	query = "select min(datechantier), max(datechantier) from chaufer"
	err = db.QueryRow(query).Scan(&first3, &last3)
	if err != nil {
		ok3 = false
	}
	//
	if !ok1 && !ok2 && !ok3 {
		return res, false, nil
	}
	if isBefore(first1, first2) && isBefore(first1, first3) {
		first = first1
	} else if isBefore(first2, first1) && isBefore(first2, first3) {
		first = first2
	} else if isBefore(first3, first1) && isBefore(first3, first2) {
		first = first3
	}
	if last1.After(last2) && last1.After(last3) {
		last = last1
	} else if last2.After(last1) && last2.After(last3) {
		last = last2
	} else if last3.After(last1) && last3.After(last2) {
		last = last3
	}
	//
	// jLim, mLim = limites de saison (jour et mois), stockées en conf
	//
	limits := strings.Split(limiteSaison, "/")
	jLim, mLim := limits[0], limits[1]
	//
	// start, end = dates de début des premières et dernières saisons
	//
	var start, end, test time.Time
	var strParse string
	// ex avec limite = 01/09 (1er sept.)
	// si first = 2018-12-15 alors start = 2018-09-01
	// si first = 2018-07-15 alors start = 2017-09-01
	strParse = strconv.Itoa(first.Year()) + "-" + mLim + "-" + jLim
	test, err = time.Parse("2006-01-02", strParse)
	if err != nil {
		return res, true, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
	}
	if test.Before(first) {
		start = test
	} else {
		strParse = strconv.Itoa(first.Year()-1) + "-" + mLim + "-" + jLim
		start, err = time.Parse("2006-01-02", strParse)
		if err != nil {
			return res, true, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
		}
	}
	// ex avec limite = 01/09 (1er sept.)
	// si last = 2020-12-15 alors end = 2020-09-01
	// si last = 2020-07-15 alors end = 2019-09-01
	// (car il s'agit de la date de début de la dernière saison)
	strParse = strconv.Itoa(last.Year()) + "-" + mLim + "-" + jLim
	test, err = time.Parse("2006-01-02", strParse)
	if err != nil {
		return res, true, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
	}
	if test.Before(last) {
		end = test
	} else {
		strParse = strconv.Itoa(last.Year()-1) + "-" + mLim + "-" + jLim
		end, err = time.Parse("2006-01-02", strParse)
		if err != nil {
			return res, true, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
		}
	}
	//
	// Calcul des dates de fin des saisons (ordre chrono inverse)
	//
	for d := end; d.After(start) || d.Equal(start); d = d.AddDate(-1, 0, 0) {
		endPeriod := d.AddDate(1, 0, -1)
		res = append(res, [2]time.Time{d, endPeriod})
	}
	return res, true, nil
}

// Auxiliaire de ComputeLimitesSaisons()
// Compare 2 dates en prenant en compte IsZero()
func isBefore(t1, t2 time.Time) bool {
	if t1.IsZero() {
		return false
	}
	if t2.IsZero() {
		return true
	}
	return t1.Before(t2)
}

func ComputeBilanValoEssences(db *sqlx.DB, dateDeb, dateFin time.Time, idsProprio []int) (valos Valorisations, err error) {
	essenceCodes := AllEssenceCodes()
	valoCodes := AllValorisationCodes() // {"PP", "CH", "PL", "PI", "BO"}
	// AllValorisationCodes() contient les valos pour les chantiers "autres valorisations"
	// Les bilans ont d'autres types de valorisation
	valoCodes = append(valoCodes, "CF") // pour séparer chauffage fermier / chauffage client
	valoCodes = append(valoCodes, "PQ") // pour ajouter plaquettes
	valos = make(Valorisations)
	for _, valoCode := range valoCodes {
		for _, essenceCode := range essenceCodes {
			valos[valoCode+"-"+essenceCode+"-vol"] = 0
			valos[valoCode+"-"+essenceCode+"-ca"] = 0
		}
	}
	//
	// chautre
	//
	chautres := []*Chautre{}
	query := `select * from chautre where datecontrat >= ? and datecontrat <= ?
	            and id in(
	                select id_chantier from chantier_ug where type_chantier='chautre' and id_ug in(
	                    select id_ug from parcelle_ug where id_parcelle in(
	                        select id from parcelle where id_proprietaire in(?)
	                    )
	                )
	            )
	`
	query, args, err := sqlx.In(query, dateDeb, dateFin, idsProprio)
	query = db.Rebind(query)
	err = db.Select(&chautres, query, args...)
	if err != nil {
		return valos, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chautre := range chautres {
		valos[chautre.TypeValo+"-"+chautre.Essence+"-vol"] += chautre.VolumeRealise
		valos[chautre.TypeValo+"-"+chautre.Essence+"-ca"] += chautre.PUHT * chautre.VolumeRealise
	}
	//
	// chaufer
	//
	chaufers := []*Chaufer{}
	query = `select * from chaufer where datechantier >= ? and datechantier <= ?
	            and id in(
	                select id_chantier from chantier_ug where type_chantier='chaufer' and id_ug in(
	                    select id_ug from parcelle_ug where id_parcelle in(
	                        select id from parcelle where id_proprietaire in(?)
	                    )
	                )
	            )
	`
	query, args, err = sqlx.In(query, dateDeb, dateFin, idsProprio)
	query = db.Rebind(query)
	err = db.Select(&chaufers, query, args...)
	if err != nil {
		return valos, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, chaufer := range chaufers {
		valos["CF-"+chaufer.Essence+"-vol"] += chaufer.Volume // CF = chauffage fermier
	}
	//
	// plaq
	//
	plaqs := []*Plaq{}
	query = `select * from plaq where datedeb >= ? and datedeb <= ?
	            and id in(
	                select id_chantier from chantier_ug where type_chantier='plaq' and id_ug in(
	                    select id_ug from parcelle_ug where id_parcelle in(
	                        select id from parcelle where id_proprietaire in(?)
	                    )
	                )
	            )
	`
	query, args, err = sqlx.In(query, dateDeb, dateFin, idsProprio)
	query = db.Rebind(query)
	err = db.Select(&plaqs, query, args...)
	if err != nil {
		return valos, werr.Wrapf(err, "Erreur query DB : "+query)
	}
	for _, plaq := range plaqs {
		err = plaq.ComputeVolume(db)
		if err != nil {
			return valos, werr.Wrapf(err, "Erreur appel plaq.computeVolume()")
		}
		valos["PQ-PS-vol"] += plaq.Volume // PQ = plaquettes ; PS = pin sylvestre
		// valos["PQ-PS-ca"] dépend des ventes ; trop compliqué (?)
	}

	//
	return valos, err
}
