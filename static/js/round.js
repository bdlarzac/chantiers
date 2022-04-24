/******************************************************************************
    Rounds a number.

    @param  x           Number to round
    @param  precision   Number of digits after decimal point.
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2020-01-30 11:15:39+01:00, Thierry Graff : Creation
********************************************************************************/
function round(x, precision){
    x = x * Math.pow(10, precision);
    x = Math.round(x);
    return (x / Math.pow(10, precision)).toFixed(precision);
}
