/******************************************************************************
   Templates

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2019-12-11 14:49:10+01:00, Thierry Graff : Creation
********************************************************************************/
package ctxt

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// used to fill Context.Template
var tmpl *template.Template

func init() {
	var fmap = template.FuncMap{
		"nl2br":             nl2br,
		"dateFr":            dateFr,
		"dateIso":           dateIso,
		"labelEssence":      labelEssence,
		"labelUnite":        labelUnite,
		"labelExploitation": labelExploitation,
		"labelValorisation": labelValorisation,
		"labelActivite":     labelActivite,
		"twoDigits":         twoDigits,
		"ucFirst":           ucFirst,
		"year":              year,
		"zero2empty":        zero2empty,
	}
	tmplDir := "view"
	tmpl = template.Must(template.New("").Funcs(fmap).ParseGlob(filepath.Join(tmplDir, "*.html"))).Option("missingkey=error")
}

// ************************* pipelines ********************************

func nl2br(t string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(t), "\n", "<br>", -1))
}

// Affiche une date JJ/MM/AAAA
func dateFr(t time.Time) template.HTML {
	return template.HTML(tiglib.DateFr(t))
}

// Affiche une date YYYY-MM-DD
func dateIso(t time.Time) template.HTML {
	return template.HTML(tiglib.DateIso(t))
}

// Affiche l'année YYYY d'une date
func year(t time.Time) template.HTML {
	return template.HTML(strconv.Itoa(t.Year()))
}

func labelEssence(str string) template.HTML {
	return template.HTML(model.LabelEssence(str))
}

// Nom des unités utilisées dans cette appli
func labelUnite(str string) template.HTML {
	return template.HTML(model.LabelUniteHTML(str))
}

// Type d'exploitation (1 - 5)
func labelExploitation(str string) template.HTML {
	return template.HTML(model.LabelExploitation(str))
}

// Type de valorisation (palette, pâte à papier...)
func labelValorisation(str string) template.HTML {
	return template.HTML(model.LabelValorisation(str))
}

// Type d'opération simple (abattage, débardage...)
func labelActivite(str string) template.HTML {
	return template.HTML(model.LabelActivite(str))
}

// from https://www.php2golang.com/method/function.ucfirst.html
func ucFirst(str string) template.HTML {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return template.HTML(u + str[len(u):])
	}
	return template.HTML("")
}

// Sert à initialiser les input type number à "" au lieu de "0"
// val doit être un int ou un float64
// Pas de vérification d'erreur
func zero2empty(val interface{}) template.HTML {
	var res string
	if _, ok := val.(float64); ok {
		if val.(float64) == 0 {
			res = ""
		} else {
			// Attention, on limite à 2 décimales - OK pour BDL
			res = strconv.FormatFloat(val.(float64), 'f', 2, 64)
		}
	} else if _, ok := val.(int); ok {
		if val.(int) == 0 {
			res = ""
		} else {
			res = strconv.Itoa(val.(int))
		}
	}
	return template.HTML(res)
}

// Pour afficher les prix, 2 chiffres après la virgule
// Des zéros sont rajoutés si besoin - par ex renvoie 12.50 au lieu de 12.5
func twoDigits(f float64) template.HTML {
	return template.HTML(fmt.Sprintf("%.2f", f))
}
