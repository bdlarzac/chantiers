/******************************************************************************
    Structure de donnée recevant le contenu de config.yml
    Pourrait aussi être dans le package ctxt.
    Mis dans model pour éviter un problème de dépendance circulaire
    lorsqu'elle est utilisée dans model.

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    @history    2019-09-26, Thierry Graff : Creation
********************************************************************************/
package model

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Paths struct {
		LogicielFoncier string `yaml:"logiciel-foncier"`
	} `yaml:"paths"`
	PourcentagePerte float64   `yaml:"pourcentage-perte"`
	DebutSaison      string    `yaml:"debut-saison"`
	TVAExt           []float64 `yaml:"tva-ext"`
	TVABDL           struct {
		Livraison           float64   `yaml:"livraison"`
		VentePlaquettes     float64   `yaml:"vente-plaquettes"`
		AutresValorisations []float64 `yaml:"autres-valorisations"`
	} `yaml:"tva-bdl"`
	Facture struct {
		// metadata - pas affiché
		Auteur   string `yaml:"auteur"`
		Createur string `yaml:"createur"`
		// Infos affichées sur les factures
		Adresse string `yaml:"adresse"`
		Tel     string `yaml:"tel"`
		Email   string `yaml:"email"`
		SiteWeb string `yaml:"site-web"`
		Siret   string `yaml:"siret"`
		TVA     string `yaml:"tva"`
	} `yaml:"facture"`
	Affacture struct {
		Adresse string `yaml:"adresse"`
	} `yaml:"affacture"`
	NbRecent int `yaml:"nb-recent"`
	Dev      struct {
		SCTLData    string `yaml:"sctl-data"`
		SCTLAnalyse string `yaml:"sctl-analyse"`
	} `yaml:"dev"`
}

// configuration spécifique au déploiement
// en utilisant les variables d'environnement
type serverEnv struct {
	DATABASE_URL string
	SCHEMA       string
	PORT         string
	MODE         string
	BACKUP_DIR   string
	CMD_PGDUMP   string
	SERVER_ADDR  string // http.Server
}

var SERVER_ENV serverEnv

// LoadEnv charge les variables d'environnement spécifiques au serveur
func MustLoadEnv() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Printf("Chargement env : %s\n", err)
	}
	SERVER_ENV = serverEnv{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		SCHEMA:       os.Getenv("SCHEMA"),
		PORT:         os.Getenv("PORT"),
		MODE:         os.Getenv("MODE"),
		CMD_PGDUMP:   os.Getenv("CMD_PGDUMP"),
		SERVER_ADDR:  os.Getenv("SERVER_ADDR"),
		BACKUP_DIR:   os.Getenv("BACKUP_DIR"),
	}
}
