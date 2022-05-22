/******************************************************************************
    Contexte, contient les variables utiles à toutes les pages

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.

    @history    2019-12-11 14:42:24+01:00, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"html/template"

	"bdl.local/bdl/model"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Page         *Page
	Redirect     string
	TemplateName string
	Template     *template.Template
	DB           *sqlx.DB
	Config       *model.Config
}

func NewContext() *Context {
	ctx := &Context{}
	ctx.Template = tmpl // déclaré dans template.go
	ctx.DB = db         // déclaré dans db.go
	ctx.Config = config // déclaré dans config.go
	return ctx
}
