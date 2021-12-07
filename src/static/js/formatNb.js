/******************************************************************************
    Remplace toLocaleString() - (pas réussi à faire marcher).
    
    @license    GPL
    @history    2021-12-06 18:52:34+01:00, Thierry Graff : Creation
********************************************************************************/

/** 
    Formate un nombre en ajoutant des espaces tous les 3 digits.
    Ex: 
        formatNb(2000) = "2 000"
        formatNb(2000.12) = "2 000.12"
        formatNb(1232000.12) = "1 232 000.12"
        formatNb(1232000) = "1 232 000"
        formatNb(20.12) = "20.12"
        formatNb(20) = "20"
    @param  x Nombre à formater
    @param  sep Caractère séparant les parties entière et décimale du nombre.
**/
function formatNb(x, sep='.'){
    [intPart, decPart] = x.toString().split(sep);
    let res = intPart[0];
    for(let i=1; i < intPart.length ; i++){
        if((intPart.length - i)%3 == 0){
            res = res + ' ';
        }
        res += intPart[i];
    }
    if(decPart != undefined){
        res += sep + decPart;
    }
    return res;
}
