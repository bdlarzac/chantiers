/**
    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
**/
package control

import (
	"archive/zip"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/generic/wilk/werr"
	"bdl.local/bdl/model"
)

type detailsAccueil struct {
	Recents []*model.Recent
}

func Accueil(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	ctx.TemplateName = "accueil.html"
	recents, err := model.GetRecents(ctx.DB)
	if err != nil {
		return err
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
	// PGPASSWORD='my_password' pg_dump --file my_dump_file -h _my_host -n my_schema -p my_port -U my_user  my_database
	cmd := exec.Command(
		model.SERVER_ENV.CMD_PGDUMP,
		model.SERVER_ENV.DATABASE_URL,
		"--file", filepath,
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
		return err
	}
	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()
	f, err := zipwriter.Create(filename)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	//
	// 3 - clean
	//
	err = os.Remove(filepath)
	if err != nil {
		return err
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
