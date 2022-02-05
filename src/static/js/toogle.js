/******************************************************************************
    Shows / hides an element.

    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2020-01-30 11:15:39+01:00, Thierry Graff : Creation
********************************************************************************/
function toogle(id){
  const elt = document.getElementById(id);
  elt.style.display = (elt.style.display == "block" ? "none" : "block");
}