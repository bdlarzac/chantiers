package control

import (
	"archive/zip"
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
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
	dirname := ctx.Config.Database.Backup.Directory
	filename := "bdl-" + time.Now().Format("2006-01-02-150405") + ".pgdump"
	filepath := dirname + string(os.PathSeparator) + filename
	//
	// 1 - pg_dump
	//
	// PGPASSWORD='my_password' pg_dump --file my_dump_file -h _my_host -n my_schema -p my_port -U my_user  my_database
	cmd := exec.Command(
		ctx.Config.Database.Backup.CmdPgdump,
		"--file", filepath,
		"-h", ctx.Config.Database.Host,
		"-p", ctx.Config.Database.Port,
		"-U", ctx.Config.Database.User,
		"-n", ctx.Config.Database.Schema,
		ctx.Config.Database.DbName)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+ctx.Config.Database.Password)
	err = cmd.Run()
	if err != nil {
		return err
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
	ctx.TemplateName = "backup.html"
	ctx.Page = &ctxt.Page{
		Header: ctxt.Header{
			Title: "Sauvegarde des données",
		},
		Menu:    "accueil",
		Details: zipfilename,
	}
	return nil
}
