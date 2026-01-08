/*
Rôles des acteurs

La notion de rôle n'est pas typée en base (pas de type role pour faire une enum).
Les strings stockées dans acteur_role.code_role correspondent aux clés de RoleMap.

@copyright  BDL, Bois du Larzac.
@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
@history    2023-04-04 16:49:32+02:00, Thierry Graff : Creation
*/
package model

// Association code rôle => label
// Les codes correspondent aux valeurs stockées en base dans acteur_role.code_role
// Attention, si on modifie les codes, il faut aussi modifier à la main
// - view/acteur-form.html
// - view/affacture-form.html
// - ComputeItems() dans model.affacture.go
// Si on ajoute une opération simple, aussi modifier
// - dans model.plaq.go, modifier la struct CoutPlaq et ComputeCouts()
var RoleMap = map[string]string{
	// Chantier plaquettes, opérations simples :
	"PLA-AB": "Abatteur", // = bûcheron
	"PLA-DB": "Débardeur",
	"PLA-DC": "Déchiqueteur",
	"PLA-BR": "Broyeur",
	"PLA-LO": "Aide logistique",
	// Chantier plaquettes, transport :
	"PLT-TR": "Transporteur PF",
	"PLT-CO": "Conducteur transport PF",
	"PLT-PO": "Propriétaire outil transport PF",
	// Chantier plaquettes, rangement :
	"PLR-RG": "Rangeur PF",
	"PLR-CO": "Conducteur rangement PF",
	"PLR-PO": "Propriétaire outil rangement PF",
	// Vente plaquettes :
	"VPL-CL": "Client PF",
	// Vente plaquettes, chargement :
	"VPC-CH": "Chargeur PF",
	"VPC-CO": "Conducteur chargement PF",
	"VPC-PO": "Propriétaire outil chargement PF",
	// Vente plaquettes, livraison
	"VPL-LI": "Livreur PF",
	"VPL-CO": "Conducteur livraison PF",
	"VPL-PO": "Propriétaire outil livraison PF",
	// Chantier autres valorisations
	"AVC-PP": "Client pâte à papier",
	"AVC-CH": "Client bois de chauffage",
	"AVC-PL": "Client palettes",
	"AVC-PI": "Client piquets",
	"AVC-BO": "Client bois d'oeuvre",
	// Divers:
	"DIV-MH": "Mesureur d'humidité",
	"DIV-PF": "Propriétaire foncier",
	"DIV-FO": "Fournisseur de plaquettes",
	"FER-BC": "Fermier bois de chauffage",
}
