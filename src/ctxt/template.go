/*
Fonctions utilisée dans les templates HTML.

	@copyright  Les fonctions spécifiques au programme BDL sont la propriété intellectuelle de BDL, Bois du Larzac.
	            Les fonctions génériques sont la propriété intellectuelle de Thierry Graff.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2019-12-11 14:49:10+01:00, Thierry Graff : Creation
*/
package ctxt

import (
	"bdl.local/bdl/generic/tiglib"
	"bdl.local/bdl/model"
	"bdl.local/bdl/view"
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// used to fill Context.Template
var tmpl *template.Template

func MustInitTemplates() {
	var fmap = template.FuncMap{
		// Generic pipelines
		"dateFr":     dateFr,
		"dateIso":    dateIso,
		"modulo":     modulo,
		"nl2br":      nl2br,
		"plus":       plus,
		"round":      round,
		"safeHTML":   safeHTML,
		"twoDigits":  twoDigits,
		"ucFirst":    ucFirst,
		"year":       year,
		"zero2empty": zero2empty,
		// Pipelines related to current program
		"labelActivite":     labelActivite,
		"labelEssence":      labelEssence,
		"labelExploitation": labelExploitation,
		"labelRole":         labelRole,
		"labelStockFrais":   labelStockFrais,
		"labelTypeVente":    labelTypeVente,
		"labelTypo":         labelTypo,
		"labelTypo_long":    labelTypo_long,
		"labelUnite":        labelUnite,
		"labelValo":         labelValo,
		"sortableUGCode":    sortableUGCode,
		"valo2uniteLabel":   valo2uniteLabel,
	}
	tmpl = template.
		Must(template.
			New("").
			Funcs(fmap).
			Option("missingkey=error").
			ParseFS(view.TemplatesFiles, "*.html", "common/*.html"))
}

// ************************* Generic pipelines ********************************

/*
Displays a date, format DD/MM/YYYY.
@copyright  Thierry Graff
@license    GPL
*/
func dateFr(t time.Time) template.HTML {
	return template.HTML(tiglib.DateFr(t))
}

/*
Displays a date, format YYYY-MM-DD.
@copyright  Thierry Graff
@license    GPL
*/
func dateIso(t time.Time) template.HTML {
	return template.HTML(tiglib.DateIso(t))
}

/*
@copyright  Thierry Graff
@license    GPL
*/
func modulo(i, mod int) int {
	return i % mod
}

/*
@copyright  Thierry Graff
@license    GPL
*/
func nl2br(t string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(t), "\n", "<br>", -1))
}

/*
@copyright  Thierry Graff
@license    GPL
*/
func plus(a, b int) int {
	return a + b
}

func round(x float64, precision int) float64 {
	x = x * math.Pow10(precision)
	x = math.Round(x)
	return x / math.Pow10(precision)
}

/*
From https://www.php2golang.com/method/function.ucfirst.html
@copyright  Thierry Graff
@license    GPL
*/
func ucFirst(str string) template.HTML {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return template.HTML(u + str[len(u):])
	}
	return template.HTML("")
}

/*
Displays the year of a date, format YYYY.
@copyright  Thierry Graff
@license    GPL
*/
func year(t time.Time) template.HTML {
	return template.HTML(strconv.Itoa(t.Year()))
}

/*
Used to initialize input type=number with "" instead of "0".
No error check.
@param      val  Must be an int or a float64.
@copyright  Thierry Graff
@license    GPL
*/
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

/*
To display prices, with a precision of 1E-2. Zeroes are added if needed.
Ex: twoDigits(12.5) returns 12.50 instead of 12.5
@copyright  Thierry Graff
@license    GPL
*/
func twoDigits(f float64) template.HTML {
	return template.HTML(fmt.Sprintf("%.2f", f))
}

/*
@copyright  Thierry Graff
@license    GPL
*/
func safeHTML(str string) template.HTML {
	return template.HTML(str)
}

// ************************* Pipelines specific to current program ********************************

// Type d'opération simple (abattage, débardage...) à partir de son code
func labelActivite(code string) template.HTML {
	return template.HTML(model.LabelActivite(code))
}

// Nom d'une essence (chêne etc.) à partir de son code
func labelEssence(code string) template.HTML {
	return template.HTML(model.EssenceMap[code])
}

// Nom d'un type d'exploitation (1 - 5), à partir de son code
func labelExploitation(code string) template.HTML {
	return template.HTML(model.LabelExploitation(code))
}

// Nom d'un type de frais pour stockage (loyer, assurance, élec) à partir de son code
func labelStockFrais(code string) template.HTML {
	return template.HTML(model.StockFraisMap[code])
}

// Nom d'un rôle (pour les acteurs) à partir de son code
func labelRole(code string) template.HTML {
	return template.HTML(model.RoleMap[code])
}

// Nom d'un type de vente (pour chautre: bois sur pied, bord de route...), à partir de son code
func labelTypeVente(code string) template.HTML {
	return template.HTML(model.ChautreTypeVenteMap[code])
}

// Nom d'une typo (couche typologique venant du PSG) utilisée dans cette appli, à partir de son code
func labelTypo(code string) template.HTML {
	return template.HTML(model.TypoMap[code])
}

// Nom d'une typo (couche typologique venant du PSG) utilisée dans cette appli, à partir de son code
// Nom complet
func labelTypo_long(code string) template.HTML {
	return template.HTML(model.TypoMap_long[code])
}

// Nom d'une unité utilisée dans cette appli, à partir de son code
func labelUnite(code string) template.HTML {
	return template.HTML(model.UniteMap[code])
}

// Type de valorisation (palette, pâte à papier...), à partir de son code
func labelValo(code string) template.HTML {
	return template.HTML(model.ValoMap[code])
}

func sortableUGCode(code string) template.HTML {
	return template.HTML(model.SortableUGCode(code))
}

// Label de l'unité correspondant à un type de valorisation (palette, pâte à papier...)
func valo2uniteLabel(code string) template.HTML {
	return template.HTML(model.UniteMap[model.CodeValo2CodeUnite(code)])
}
