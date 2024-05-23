/** 
    Code permettant de modifier la date de paiement directement dans un liste.
    Utilisé par 
        - venteplaq-list.html
        - chautre-list.html
    La template utilisatrice doit :
        - Avoir une <div id="message-modif-date-paiement"></div>
        
    "entity" désigne la chose à modifier = vente plaquette ou chantier autre valorisation.
    
    @param  urlAjax     url à appeler pour faire la modif en base via ajax.
                        = "/ajax/update/date-venteplaq" ou "/ajax/update/chautre"
    @param  idEntity    Id en base de la vente ou du chantier à modifier
    @param  titreEntity Ne sert que pour les messages ; titre de la vente ou du chantier qui est modifié.e
    @param  labelEntity Ne sert que pour les messages ; = "la vente" ou "le chantier"
**/
async function updateDatePaiement(urlAjax, idEntity, titreEntity, labelEntity){
    // background du message affiché suite à la modif
    const rootVars = getComputedStyle(document.querySelector(':root'));
    const bgOk =  rootVars.getPropertyValue('--background-ok');
    const bgError =  rootVars.getPropertyValue('--background-error');
    
    let msg;
    let ok; // 'ok' ou 'nok'
    const dateElt = document.getElementById('update-paiement-' + idEntity);
    let newDate = dateElt.value;
    const DATE_NULLE = '0001-01-01'; // date nulle pour postgres - TODO ne devrait pas être en dur comme ça
    if(newDate == ""){
        // Se produit lorsque l'utilisateur modifie une date renseignée
        // pour la passer à non renseignée (clic sur le bouton "clear" du calendrier)
        newDate = DATE_NULLE;
    }
    const url = urlAjax + '/' + idEntity + '/' + newDate;
    const response = await fetch(url);
    if(response == null){
        msg = "- ERREUR - Transmettez ce message à l'administrateur du site :\n"
            + "Problème de récupération " + url;
        ok = 'nok';
    }
    else{
        result = await response.json();
        ok = result['ok']; 
        if(ok == 'ok'){
            if(newDate != DATE_NULLE){
                msg = `<b>Modification enregistrée</b> pour ${labelEntity} ${idEntity} '${titreEntity}'`;
            }
            else {
                msg = `<b>Date de paiement effacée</b> pour ${labelEntity} ${idEntity} '${titreEntity}'`;
            }
        }
        else{
            msg = "<b>PROBLEME</b>, modification non effectuée : \n"
                + result['message'];
        }
    }
    const msgElt = document.getElementById('message-modif-date-paiement');
    msgElt.innerHTML = msg;
    msgElt.style.display='inline-block';
    let timeout;
    if(ok == 'ok'){
        msgElt.style.background = bgOk;
        if(newDate != DATE_NULLE){
            dateElt.style.background = bgOk;
        }
        else {
            dateElt.style.background = 'white'; // bof, le background par défaut devrait pouvoir être calculé
        }
        timeout = 3000;
    }
    else {
        msgElt.style.background = bgError;
        dateElt.style.background = bgError;
        timeout = 6000;
    }
    // faire disparaître le message au bout d'un moment
    setTimeout(function(){
        document.getElementById('message-modif-date-paiement').style.display = 'none';
    }, timeout);
                     
}