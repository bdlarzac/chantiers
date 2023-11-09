/*
Code liés aux saisons

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2021-01-19 10:09:42+01:00, Thierry Graff : Creation
@history    2023-05-17 09:16:48+02:00, Thierry Graff : Refactor, isole le code lié aux saisonx
*/
package model

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

// Renvoie un tableau contenant les dates de début / fin des "saisons"
// Les saisons encadrent tous les chantiers et ventes stockés en base.
// Une saison dure un an.
//
// @param limiteSaison string au format JJ/MM (tiré de 'debut-saison' en conf)
// @return
//   - un tableau de 2 time.Time avec les dates limites des saisons
//   - un bool indiquant s'il existe des chantiers ou des ventes en base
//   - une erreur éventuelle
func ComputeLimitesSaisons(db *sqlx.DB, limiteSaison string) ([][2]time.Time, bool, error) {
	// retour
	var res [][2]time.Time
	var err error
	//
	// first, last = dates du premier et dernier chantier ou vente en base
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
	// ventes plaquettes
	ok4 := true
	var first4, last4 time.Time
	query = "select min(datevente), max(datevente) from venteplaq"
	err = db.QueryRow(query).Scan(&first4, &last4)
	if err != nil {
		ok4 = false
	}
	//
	if !ok1 && !ok2 && !ok3 && !ok4 {
		return res, false, nil
	}
	if tiglib.IsBefore(first1, first2) && tiglib.IsBefore(first1, first3) && tiglib.IsBefore(first1, first4) {
		first = first1
	} else if tiglib.IsBefore(first2, first1) && tiglib.IsBefore(first2, first3) && tiglib.IsBefore(first2, first4) {
		first = first2
	} else if tiglib.IsBefore(first3, first1) && tiglib.IsBefore(first3, first2) && tiglib.IsBefore(first3, first4) {
		first = first3
	} else if tiglib.IsBefore(first4, first1) && tiglib.IsBefore(first4, first2) && tiglib.IsBefore(first4, first3) {
		first = first4
	}
	if last1.After(last2) && last1.After(last3) && last1.After(last4) {
		last = last1
	} else if last2.After(last1) && last2.After(last3) && last2.After(last4) {
		last = last2
	} else if last3.After(last1) && last3.After(last2) && last3.After(last4) {
		last = last3
	} else if last4.After(last1) && last4.After(last2) && last4.After(last3) {
		last = last4
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
