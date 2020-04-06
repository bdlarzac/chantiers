/******************************************************************************
    Arrondit un nombre à n chiffres après la virgule.

    @license    GPL
    @history    2020-01-30 11:15:39+01:00, Thierry Graff : Creation
********************************************************************************/


// *****************************************
/** 
    @param  x Nombre à arrondir
    @param  precision Nombre de chiffres après la virgule
**/
function round(x, precision){
    x = x * Math.pow(10, precision);
    x = Math.round(x);
    return (x / Math.pow(10, precision)).toFixed(precision);
    
}
