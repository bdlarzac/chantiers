package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"net/http"
	"time"
	"os"
	"os/exec"
    "archive/zip"
    "io/ioutil"
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
    // PGPASSWORD='bdl2019-09-29' pg_dump --file bdl.pgsql.dump -h localhost -n public -p 5432 -U bdl  bdlbois
    var err error
    dirname := ctx.Config.Database.Backup.Directory
    filename := "bdl-" + time.Now().Format(time.RFC3339)[:19] + ".pgdump"
    filepath := dirname + string(os.PathSeparator) + filename
    //
    // 1 - pg_dump
    //
    cmd := exec.Command(
        ctx.Config.Database.Backup.CmdPgdump,
        "--file", filepath,
        "-h", ctx.Config.Database.Host,
        "-p", ctx.Config.Database.Port,
        "-U", ctx.Config.Database.User,
        "-n", ctx.Config.Database.Schema,
        ctx.Config.Database.DbName)
	cmd.Env = append(os.Environ(), "PGPASSWORD=" + ctx.Config.Database.Password)
    err = cmd.Run()
    if err != nil {
        return err
    }
    //
    // 2 - zip
    //
    zipfilename := filepath + ".zip"
    zipfile, err := os.Create(zipfilename)
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
		Menu: "accueil",
		Details: zipfilename,
	}
	return nil
}
