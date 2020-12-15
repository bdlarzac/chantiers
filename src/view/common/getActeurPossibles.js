/** 
    Promesse utilisée par autocomplete() pour s'alimenter en données
    Utilisé par les formulaires avec un autocomplete sur nom acteur
    
    ATTENTION : remmplit une variable globale window.globalPersons
        (utilisée dans la validation des formulaires pour s'assurer que le nom
        de la personne n'a pas été modifié entre le retour de l'autocomplete
        et la validation du formulaire)
    
    @param inputField input type=text utilisé pour récupérer des noms de personnes par ajax
**/
async function getActeursPossibles(inputField){
    let result = [];
    const url = "/ajax/autocomplete/acteur/" + inputField.value;
    let response;
    response = await fetch(url);
    if(response == null){
        alert("- ERREUR - Transmettez ce message l'administrateur du site :\n"
            + "Problème de récupération " + url);
        return result;
    }
    response = await response.json();
    response.forEach(function(item) {
        result.push(item.name);
        window.globalPersons[item.id] = item.name;
    });
    return result;
}
