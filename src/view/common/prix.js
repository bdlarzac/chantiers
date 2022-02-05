/******************************************************************************
    Opérations sur les prix.

    @copyright  Thierry Graff.
    @licence    GPL, conformémént au fichier LICENCE situé à la racine du projet.
    
    @history    2020-02-10 22:11:54+01:00, Thierry Graff : Creation
********************************************************************************/


/** 
    Calcule prix HT
    @param  pu  Prix unitaire
    @param  qte Quantité
**/
function prixHT(pu, qte){
    return pu * qte;
}

/** 
    Calcule prix TTC
    @param  ht  Prix HT
    @param  tva Taux de TVA
**/
function prixTTC(ht, tva){
    return ht * (1 + (tva/100));
}
