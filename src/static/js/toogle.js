/******************************************************************************
    Cache / montre un élément.

    @license    GPL
    @history    2020-01-30 11:15:39+01:00, Thierry Graff : Creation
********************************************************************************/

function toogle(id){
  var elt = document.getElementById(id);
  elt.style.display == "block" ? elt.style.display = "none" : elt.style.display = "block"; 
}