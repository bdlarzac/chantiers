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
		"dateFr":            dateFr,
		"dateIso":           dateIso,
		"labelEssence":      labelEssence,
		"labelUnite":        labelUnite,
		"labelExploitation": labelExploitation,
		"labelValorisation": labelValorisation,
		"labelActivite":     labelActivite,
		"labelStockFrais":   labelStockFrais,
		"nl2br":             nl2br,
		"safeHTML":          safeHTML,
		"twoDigits":         twoDigits,
		"ucFirst":           ucFirst,
		"year":              year,
		"zero2empty":        zero2empty,
	}
	tmpl = template.
		Must(template.
			New("").
			Funcs(fmap).
			ParseGlob(filepath.Join("view", "*.html"))).
		Option("missingkey=error")
	tmpl.New("chantier-lien").Funcs(fmap).ParseFiles(filepath.Join("view", "common", "chantier-lien.html"))
	tmpl.New("chantier-lien-help").Funcs(fmap).ParseFiles(filepath.Join("view", "common", "chantier-lien-help.html"))
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

// Nom d'une essence (chêne etc.) à partir de son code
func labelEssence(str string) template.HTML {
	return template.HTML(model.LabelEssence(str))
}

// Nom d'une unité utilisée dans cette appli, à partir de son code
func labelUnite(str string) template.HTML {
	return template.HTML(model.LabelUniteHTML(str))
}

// Type d'exploitation (1 - 5), à partir de son code
func labelExploitation(str string) template.HTML {
	return template.HTML(model.LabelExploitation(str))
}

// Type de valorisation (palette, pâte à papier...), à partir de son code
func labelValorisation(str string) template.HTML {
	return template.HTML(model.LabelValorisation(str))
}

// Type d'opération simple (abattage, débardage...) à partir de son code
func labelActivite(str string) template.HTML {
	return template.HTML(model.LabelActivite(str))
}

// Type de frais pour stockage (loyer, assurance, élec) à partir de son code
func labelStockFrais(str string) template.HTML {
	return template.HTML(model.LabelStockFrais(str))
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

func safeHTML(str string) template.HTML {
	return template.HTML(str)
}
