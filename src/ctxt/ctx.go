/******************************************************************************
    Contexte, contient les variables utiles à toutes les pages

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-11 14:42:24+01:00, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"bdl.local/bdl/model"
	"github.com/jmoiron/sqlx"
	"html/template"
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
	ctx.Template = tmpl
	ctx.DB = db
	ctx.Config = config
	return ctx
}
