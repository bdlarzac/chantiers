/**
    Génère une liste d'options html à partir d'un slice
    Adaptation d'un code de wilk

    Changements par rapport au code de wilk :
    Utilisé que FmtOption
    - déplace l'espace pour afficher selected
    - utilise " au lieu de '
    - supprime des sauts de ligne
    - Ajoute OptionId dans OptionString et OptionInt
      (il faudrait faire pareil pour OptionGroupString et OptionGroupInt)
      Motivé par l'utilisation de plusieurs selects avec les mêmes options dans la même page (taux TVA)
      => pour que chaque option ait un id unique
      Comportement : si OptionId est fourni, l'id de l'option vaut OptionId, sinon vaut OptionValue
                     Dans tous les cas les options ont un id
**/

package webo

import (
	"bytes"
	"html/template"
)

var tmpOptions = template.Must(template.New("options").Parse(`{{$sel := .Sel}}{{range .Options}}<option {{if eq .OptionId ""}}id="{{.OptionValue}}"{{else}}id="{{.OptionId}}"{{end}} value="{{.OptionValue}}"{{if eq .OptionValue $sel}} selected{{end}}>{{.OptionLabel}}</option>
{{end}}`)).Option("missingkey=error")

type OptionString struct {
	OptionValue string
	OptionId    string
	OptionLabel string
}
type OptionInt struct {
	OptionValue int
	OptionId    string
	OptionLabel string
}

func FmtOptions(slice interface{}, sel interface{}) template.HTML {
	var s bytes.Buffer
	err := tmpOptions.Execute(&s, map[string]interface{}{"Options": slice, "Sel": sel})
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(s.String())
}

var tmpOptionsGroup = template.Must(template.New("optionsGroup").Parse(`
{{$sel := .Sel}}
{{$group := ""}}
{{range .Options}}
	{{if ne $group .OptionGroup}}
	{{if ne $group ""}} </optgroup> {{end}}
	<optgroup label='{{.OptionGroup}}'>
	{{$group = .OptionGroup}}
	{{end}}
<option value='{{.OptionValue}}' {{if eq .OptionValue $sel}}selected{{end}}>{{.OptionLabel}}</option>
{{end}}
</optgroup>
`)).Option("missingkey=error")

type OptionGroupString struct {
	OptionGroup string
	OptionValue string
	OptionLabel string
}
type OptionGroupInt struct {
	OptionGroup string
	OptionValue int
	OptionLabel string
}

func FmtOptionsGroup(slice interface{}, sel interface{}) template.HTML {
	var s bytes.Buffer
	err := tmpOptionsGroup.Execute(&s, map[string]interface{}{"Options": slice, "Sel": sel})
	if err != nil {
		return template.HTML(err.Error())
	}
	return template.HTML(s.String())
}
