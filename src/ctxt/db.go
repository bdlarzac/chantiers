/******************************************************************************

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-09-26 23:41:06+02:00, Thierry Graff : Creation
********************************************************************************/

package ctxt

import (
	"log"
	"strings"

	"bdl.local/bdl/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// ajoute le schema a l'url
// et supprime sslmode=prefer mis automatiquement par Scalingo
func AjusteDbURL(url string, schema string) string {
	dbURL := strings.ReplaceAll(url, "sslmode=prefer", "") // non supporté par lib/pq
	if strings.HasPrefix(dbURL, "postgres") && !strings.Contains(dbURL, "search_path") {
		if !strings.Contains(dbURL, "?") {
			dbURL += "?"
		} else {
			dbURL += "&"
		}
		dbURL += "search_path=" + schema
//		dbURL += "&sslmode=disable"
	}
	return dbURL
}

func MustInitDB() {
	var err error
//	dbURL := AjusteDbURL(model.SERVER_ENV.DATABASE_URL, model.SERVER_ENV.DATABASE_SCHEMA)
//	db, err = sqlx.Open("postgres", dbURL)
	db, err = sqlx.Open(
	    "postgres",
	    "dbname="+model.SERVER_ENV.DATABASE_DBNAME+
	    " user="+model.SERVER_ENV.DATABASE_USER+
	    " password="+model.SERVER_ENV.DATABASE_PASSWORD+
	    " host="+model.SERVER_ENV.DATABASE_HOST+
	    " port="+model.SERVER_ENV.DATABASE_PORT+
	    " search_path="+model.SERVER_ENV.DATABASE_SCHEMA+
	    " sslmode="+model.SERVER_ENV.DATABASE_SSLMODE,
	)
	if err != nil {
		log.Fatalf("Connexion DB impossible : %v", err)
	}
	//TODO: ici faire upgrade versions
}
