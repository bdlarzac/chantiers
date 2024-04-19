/******************************************************************************
    Code relatif aux ventes de plaquettes
    Mis ici car commun à plusieurs pages

    @copyright  BDL, Bois du Larzac.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    
    @history    2020-03-20 10:59:03+01:00, Thierry Graff : Creation
********************************************************************************/

function deleteVentePlaquette(idVente, nomClient, dateVente){
    let msg = "En cliquant sur OK,"
            + "\nla vente \"" + nomClient + " " + dateVente + "\" sera définitivement supprimée."
            + "\n\nATTENTION, cette action a les conséquences suivantes :"
            + "\n- Supprime toutes les livraisons associées à cette vente."
            + "\n- Supprime tous les chargements associés à ces livraisons."
            + "\n- Rétablit les stocks des tas associés à ces chargements.\n";
    let r = confirm(msg);
    if (r == true) {
        window.location = "/vente/delete/" + idVente;
    }
}


/** 
    Sert de filtre pour afficher uniquement les factures contenant suffisamment d'information.
**/
function showVentePlaqFacture(idVente, numFacture, dateFacture) {
    if(numFacture == "" || dateFacture == ""){
        alert(
            "Pour voir la facture, il faut d'abord renseigner :"
            + "\n- le numéro de facture"
            + "\n- la date de facture"
            + "\nModifiez la vente et renseignez ces 2 champs."
        );
        return;
    }
    window.location = "/facture/vente-plaquette/" + idVente;
}
