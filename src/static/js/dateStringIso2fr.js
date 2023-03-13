/******************************************************************************

    Converts a string YYYY-MM-DD to DD/MM/YYYY
    WARNING: does not check input parameter
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2023-03-13 10:02:57+01:00, Thierry Graff : Creation
********************************************************************************/
function dateStringIso2fr(yyyy_mm_dd){
    const y = yyyy_mm_dd.substring(0, 4);
    const m = yyyy_mm_dd.substring(5, 7);
    const d = yyyy_mm_dd.substring(8);
    return d + '/' + m + '/' + y;
}
