/******************************************************************************
    Calculs utilisés dans les bilans

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2021-01-19 10:09:42+01:00, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

// Renvoie un tableau contenant les dates de début / fin des "saisons"
// Les saisons encadrent tous les chantiers plaquettes stockés en base.
// Une saison dure un an
// @param limiteSaison string au format JJ/MM (tiré de 'debut-saison' en conf)
func ComputeLimitesSaisons(db *sqlx.DB, limiteSaison string) ([][2]time.Time, error) {
	var res [][2]time.Time
	var err error
	//
	// first, last = date du premier et dernier chantier plaquettes en base
	//
	var first, last time.Time
	var query string
	query = "select min(datedeb), max(datedeb) from plaq"
	err = db.QueryRow(query).Scan(&first, &last)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur query : "+query)
	}
	//
	// dLim, mLim = limites de saison, stockées en conf
	//
	limits := strings.Split(limiteSaison, "/")
	dLim, mLim := limits[0], limits[1]
	//
	// start, end = dates de début des premières et dernières saisons
	//
	var start, end, test time.Time
	var strParse string
	// ex avec limite = 01/09 (1er sept.)
	// si first = 2018-12-15 alors start = 2018-09-01
	// si first = 2018-07-15 alors start = 2017-09-01
	strParse = strconv.Itoa(first.Year()) + "-" + mLim + "-" + dLim
	test, err = time.Parse("2006-01-02", strParse)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
	}
	if test.Before(first) {
		start = test
	} else {
		strParse = strconv.Itoa(first.Year()-1) + "-" + mLim + "-" + dLim
		start, err = time.Parse("2006-01-02", strParse)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
		}
	}
	// ex avec limite = 01/09 (1er sept.)
	// si last = 2020-12-15 alors end = 2020-09-01
	// si last = 2020-07-15 alors end = 2019-09-01
	// (car il s'agit de la date de début de la dernière saison)
	strParse = strconv.Itoa(last.Year()) + "-" + mLim + "-" + dLim
	test, err = time.Parse("2006-01-02", strParse)
	if err != nil {
		return res, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
	}
	if test.Before(last) {
		end = test
	} else {
		strParse = strconv.Itoa(last.Year()-1) + "-" + mLim + "-" + dLim
		end, err = time.Parse("2006-01-02", strParse)
		if err != nil {
			return res, werr.Wrapf(err, "Erreur appel time.Parse("+strParse+")")
		}
	}
	//
	// Calcul des dates de fin des saisons (ordre chrono inverse)
	//
	for d := end; d.After(start) || d.Equal(start); d = d.AddDate(-1, 0, 0) {
		endPeriod := d.AddDate(1, 0, -1)
		res = append(res, [2]time.Time{d, endPeriod})
	}
	return res, nil
}
