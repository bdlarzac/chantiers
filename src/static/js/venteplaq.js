/******************************************************************************
    Code relatif aux ventes de plaquettes
    Mis ici car commun à plusieurs pages

    @license    GPL
    @history    2020-03-20 10:59:03+01:00, Thierry Graff : Creation
********************************************************************************/


// *****************************************
function deleteVentePlaquette(idVente, nomClient, dateVente){
    let msg = "ATTENTION, en cliquant sur OK,\n"
            + "la vente \"" + nomClient + " " + dateVente + "\" sera définitivement supprimée.\n"
            + "\nAttention car les livraisons et chargements"
            + "\nassociés à cette vente seront aussi supprimés.\n";
    let r = confirm(msg);
    if (r == true) {
        window.location = "/vente/delete/" + idVente;
    }
}


// *****************************************
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
