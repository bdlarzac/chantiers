/******************************************************************************

    @copyright  Thierry Graff
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-09-26 23:41:06+02:00, Thierry Graff : Creation
********************************************************************************/

package ctxt

import (
	"log"
	"bdl.local/bdl/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func MustInitDB() {
	var err error
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
