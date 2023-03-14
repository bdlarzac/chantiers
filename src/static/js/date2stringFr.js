/******************************************************************************

    Converts a Date object to a string formatted DD/MM/YYYY
    WARNING: does not check input parameter
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2023-03-13 17:28:03+01:00, Thierry Graff : Creation
********************************************************************************/
function date2stringFr(date){
    // const y = date.getYear();
    // const m = date.getMonth();
    // const d = date.getDay();
    return date.toLocaleString('en-GB').substring(0, 10);
}
