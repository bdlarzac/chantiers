/*
Gestion des chantiers visités récemment

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2020-05-14 17:35:22+02:00, Thierry Graff : Creation
*/
package model

import (
	"bdl.local/bdl/generic/wilk/werr"
	"github.com/jmoiron/sqlx"
	"time"
)

type Recent struct {
	URL        string
	Label      string
	DateVisite time.Time
}

// Renvoie les dernières activités visitées
// (les plus récentes sont renvoyées en premier)
func GetRecents(db *sqlx.DB) ([]*Recent, error) {
	r := []*Recent{}
	query := "select * from recent order by datevisite desc"
	err := db.Select(&r, query)
	if err != nil {
		return r, werr.Wrapf(err, "Erreur query : "+query)
	}
	return r, nil
}

// Modifie l'historique
// si l'url est déjà en base, update la ligne en changeant la date de visite
// sinon ajoute une ligne.
// Ensuite efface les entrées les plus anciennes
//
// Attention, l'historique basée sur les urls peut mener à des doublons :
// /chantier/plaquette/1
// /chantier/plaquette/1/chantiers
// ou
// /chantier/chauffage-fermier/liste
// /chantier/chauffage-fermier/liste/2021
// La logique de dédoublonnage n'est pas gérée ici,
// mais dans les contrôleurs, dans le code qui appelle AddRecent()
func AddRecent(db *sqlx.DB, conf *Config, r *Recent) error {
	var err error
	var count int
	now := time.Now()
	_ = db.QueryRow("select count(*) from recent where url=$1", r.URL).Scan(&count)
	if count == 0 {
		query := `insert into recent(
            url,
            label,
            datevisite
            ) values($1,$2,$3)`
		_, err = db.Exec(
			query,
			r.URL,
			r.Label,
			now)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	} else {
		// update aussi le label en cas d'update d'une activité
		query := `update recent set (datevisite,label)=($1,$2) where url=$3`
		_, err = db.Exec(query, now, r.Label, r.URL)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	// delete les plus anciens
	all := []*Recent{}
	query := "select * from recent order by datevisite desc" // plus récente visite en premier
	err = db.Select(&all, query)
	if err != nil {
		return werr.Wrapf(err, "Erreur query : "+query)
	}
	query = "delete from recent where url=$1"
	for i, r := range all {
		if i < conf.NbRecent {
			continue
		}
		_, err := db.Exec(query, r.URL)
		if err != nil {
			return werr.Wrapf(err, "Erreur query : "+query)
		}
	}
	return nil
}
