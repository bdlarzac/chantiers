/******************************************************************************
    Shows / hides an element.
    Same as toogle(), but works inside a TR (table row).
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2021-02-09 12:29:51+01:00, Thierry Graff : Creation
********************************************************************************/
function toogleTR(id){
  const elt = document.getElementById(id);
  elt.style.display = (elt.style.display == "table-row" ? "none" : "table-row");
}