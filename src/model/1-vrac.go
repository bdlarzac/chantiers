/*
*****************************************************************************

	Diverses fonctions utilitaires

	@copyright  BDL, Bois du Larzac.
	@licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
	@history    2021-03-01 15:33:27+01:00, Thierry Graff : Creation

*******************************************************************************
*/
package model

// Convertit un volume de map vert en un volume de map sec
// en prenant en compte un pourcentage de perte (stocké en config)
// Fonction mise ici car pourrait être dans plaq.go ou tas.go ou stockage.go
func Vert2sec(vert, pourcentPerte float64) float64 {
	return vert * (100 - pourcentPerte) / 100
}
