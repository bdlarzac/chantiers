/** 
    Auxiliaire des formulaires pour vérifier les input text contenant un nom d'acteur.
    ATTENTION, utilise window.globalPersons, qui doit être définie dans le formulaire
    
    @param  eltId       id de l'input type=text contenant le nom d'acteur
    @param  msgEmpty    Message à afficher si le champ est vide
    @return Tableau avec 2 éléments
                - id de l'acteur en base, ou 0 si inexistant
                - Message d'erreur à afficher, ou "" si vérification OK
**/
function checkActeur(eltId, msgEmpty){
    let msg = "", idActeur = 0;
    const personName = document.getElementById(eltId).value.trim();
    if(personName == ""){
        msg += msgEmpty;
    }
    else if(window.globalPersons.length == 0){
        // se produit lorsque le user rentre n'importe quoi
        // sans qu'un appel ajax n'ait été fait auparavant
        msg += "- Le nom \"" + personName + "\" ne correspond à aucun acteur connu en base.\n";
    }
    else{
        ok = false;
        for(idActeur in window.globalPersons){
            if(window.globalPersons[idActeur] == personName){
                ok = true;
                break;
            }
        }
        if(!ok){
            msg += "- Le nom \"" + personName + "\" ne correspond à aucun acteur connu en base.\n";
        }
    }
    return [idActeur, msg];
}
