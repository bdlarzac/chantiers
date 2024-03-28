/*
Contrôle l'affichage de la page d'accueil et des pages générales (menu accueil).

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
*/
package control

import (
	"archive/zip"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type detailsAccueil struct {
	Recents []*model.Recent
	Periods [][2]time.Time
}

func Accueil(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "accueil.html"
	recents, err := model.GetRecents(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}

	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Accueil",
		},
		Menu: "accueil",
		Details: detailsAccueil{
			Recents: recents,
		},
	}
	return nil
}

func ShowDoc(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "doc.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Documentation",
		},
		Menu: "accueil",
	}
	return nil
}

func BackupDB(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	var err error
	dirname := model.SERVER_ENV.BACKUP_DIR
	filename := "bdl-" + time.Now().Format("2006-01-02-150405") + ".pgdump"
	filepath := dirname + string(os.PathSeparator) + filename
	//
	// 1 - pg_dump
	//
	// PGPASSWORD='my_password' pg_dump --file <pg_dump_file> -h <pg_host> -n <pg_schema> -p <pg_port> -U <pg_user> <pg_database>
	// dbURL=postgresql://<pg_user>:<pg_password>@<pg_host>:<pg_port>/<pg_database>?sslmode=<pg_sslmode>
	dbURL := "postgresql://"+
	    model.SERVER_ENV.DATABASE_USER+
	    ":"+model.SERVER_ENV.DATABASE_PASSWORD+
	    "@"+model.SERVER_ENV.DATABASE_HOST+
	    ":"+model.SERVER_ENV.DATABASE_PORT+
	    "/"+model.SERVER_ENV.DATABASE_DBNAME+
	    "?sslmode="+model.SERVER_ENV.DATABASE_SSLMODE
	cmd := exec.Command(
		model.SERVER_ENV.CMD_PGDUMP,
		dbURL,
		"--file="+filepath,
		"--schema="+model.SERVER_ENV.DATABASE_SCHEMA,
	)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if model.SERVER_ENV.RUN_MODE == "prod" {
			return werr.Wrapf(err, "Paramètres BDD invalides")
		}
		return werr.Wrapf(err, "Paramètres BDD invalides:\n%v", cmd.Args)
	}
	//
	// 2 - zip
	//
	zipfilepath := filepath + ".zip"
	zipfilename := filename + ".zip"
	zipfile, err := os.Create(zipfilepath)
	if err != nil {
		return werr.Wrap(err)
	}
	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()
	f, err := zipwriter.Create(filename)
	if err != nil {
		return werr.Wrap(err)
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return werr.Wrap(err)
	}
	_, err = f.Write(b)
	if err != nil {
		return werr.Wrap(err)
	}
	//
	// 3 - clean
	//
	err = os.Remove(filepath)
	if err != nil {
		return werr.Wrap(err)
	}
	//
	// 4 - affiche
	//
	ctx.TemplateName = "db-backup.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Sauvegarde des données",
		},
		Menu:    "accueil",
		Details: zipfilename,
	}
	return nil
}

// Demande de recalculer la table qgis_export
func MajQGis(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	err := model.QGisUpdate(ctx.DB)
	if err != nil {
		return werr.Wrap(err)
	}
	//
	ctx.TemplateName = "qgis-update.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Mise à jour pour QGis",
		},
		Menu: "accueil",
	}
	return nil
}

type detailsBlocNotes struct {
	OK        bool
	Contenu   string
	UrlAction string
}

func UpdateBlocnotes(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "POST":
		//
		// Process form
		//
		if err = r.ParseForm(); err != nil {
			return werr.Wrap(err)
		}
		contenu := r.PostFormValue("contenu")
		err = model.UpdateBlocnotes(ctx.DB, contenu)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.Redirect = "/bloc-notes/update/ok"
		return nil
	default:
		//
		// Affiche form
		//
		vars := mux.Vars(r)
		_, ok := vars["ok"]
		contenu, err := model.GetBlocnotes(ctx.DB)
		if err != nil {
			return werr.Wrap(err)
		}
		ctx.TemplateName = "bloc-notes-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Modifier le bloc-notes",
				CSSFiles: []string{
					"/static/css/form.css",
				},
			},
			Menu: "accueil",
			Details: detailsBlocNotes{
				Contenu:   contenu,
				UrlAction: "/bloc-notes/update/ok",
				OK:        ok,
			},
		}
		return nil
	}
}
