/******************************************************************************
    Code relatif aux chantiers autres valorisations
    Mis ici car commun à plusieurs pages (chautre-show et chautre-list)

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    
    @history    2024-04-30 05:30:34+02:00, Thierry Graff : Creation
********************************************************************************/

/** 
    Sert de filtre pour afficher uniquement les factures contenant suffisamment d'information.
**/
function ShowFactureChautre(idChantier, numFacture, dateFacture) {
    if(numFacture == "" || dateFacture == ""){
        alert(
            "Pour voir la facture, il faut d'abord renseigner :"
            + "\n- le numéro de facture"
            + "\n- la date de facture"
            + "\nModifiez le chantier et renseignez ces 2 champs."
        );
        return;
    }
    window.location = "/facture/autre/" + idChantier;
}

