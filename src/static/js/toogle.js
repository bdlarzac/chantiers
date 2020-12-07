/******************************************************************************
    Cache / montre un élément.

    @license    GPL
    @history    2020-01-30 11:15:39+01:00, Thierry Graff : Creation
********************************************************************************/

function toogle(id){
  const elt = document.getElementById(id);
  elt.style.display = (elt.style.display == "block" ? "none" : "block");
}