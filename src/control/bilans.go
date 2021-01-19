package control

import (
	"bdl.local/bdl/ctxt"
	"bdl.local/bdl/model"
	"net/http"
	"time"
)

type detailsBilansForm struct {
    Periods [][2]time.Time
	UrlAction string
}

func FormBilans(ctx *ctxt.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		/*
		   		//
		   		// Process form
		   		//
		   		chantier, err := chantierBSPiedForm2var(r)
		   		if err != nil {
		   			return err
		   		}
		   		// calcul des ids UG, Lieudit et Fermier, pour transmettre à InsertBSPied()
		           var idsUG, idsLieudit, idsFermier []int
		   		var id int
		           for key, val := range(r.PostForm){
		               if strings.Index(key, "ug-") == 0 {
		                   // ex : ug-0:[6] (6 est l'id UG)
		                   id, err = strconv.Atoi(val[0])
		                   if err != nil {
		                       return err
		                   }
		                   idsUG = append(idsUG, id)
		               }
		               if strings.Index(key, "lieudit-") == 0 {
		                   // ex : lieudit-164:[on] (164 est l'id lieudit)
		                   id, err = strconv.Atoi(key[8:])
		                   if err != nil {
		                       return err
		                   }
		                   idsLieudit = append(idsLieudit, id)
		               }
		               if strings.Index(key, "fermier-") == 0 {
		                   // ex : fermier-25:[on] (25 est l'id fermier)
		                   id, err = strconv.Atoi(key[8:])
		                   if err != nil {
		                       return err
		                   }
		                   idsFermier = append(idsFermier, id)
		               }
		           }
		           //
		   		_, err = model.InsertBSPied(ctx.DB, chantier, idsUG, idsLieudit, idsFermier)
		   		if err != nil {
		   			return err
		   		}
		   		ctx.Redirect = "/chantier/bois-sur-pied/liste/" + strconv.Itoa(chantier.DateContrat.Year())
		   		// model.AddRecent() inutile puisqu'on est redirigé vers la liste, où AddRecent() est exécuté
		*/
		return nil
	default:
		//
		// Affiche form
		//
		periods, err := model.ComputeLimitesSaisons(ctx.DB, ctx.Config.DebutSaison)
		if err != nil {
			return err
		}
		ctx.TemplateName = "bilans-form.html"
		ctx.Page = &ctxt.Page{
			Header: ctxt.Header{
				Title: "Choix bilan",
				CSSFiles: []string{
					"/static/css/form.css"},
			},
			Details: detailsBilansForm{
				Periods:  periods,
				UrlAction: "/bilans/show",
			},
			Menu: "bilans",
			Footer: ctxt.Footer{
				JSFiles: []string{
					"/static/js/toogle.js",
				},
			},
		}
		return nil
	}
}
