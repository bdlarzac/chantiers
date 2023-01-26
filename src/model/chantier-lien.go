/******************************************************************************
    Code commune aux différents types de chantier
    Pour gérer les liens entre chantiers et autres entités

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2023-01-26 11:16:37+01:00, Thierry Graff : Creation à partir de ChauferParcelle
********************************************************************************/
package model

/** 
    Lien entre une parcelle et un chantier (table chantier_parcelle)
    Utilisé par plaq, chautre, chaufer
    Pour chaque parcelle, on doit préciser s'il s'agit d'une parcelle entière ou pas.
    S'il ne s'agit pas d'une parcelle entière, il faut préciser la surface concernée par la coupe.
**/
type ChantierParcelle struct {
	IdChantier  int `db:"id_chantier"`
	IdParcelle int `db:"id_parcelle"`
	Entiere    bool
	Surface    float64
	// Pas stocké en base
	Parcelle *Parcelle
}
