/******************************************************************************
    Cache / montre un élément.
    Même chose que toogle(), mais fonctionne à l'intérieur d'une TR (table row)
    @license    GPL
    @history    2021-02-09 12:29:51+01:00, Thierry Graff : Creation
********************************************************************************/

function toogleTR(id){
  const elt = document.getElementById(id);
  elt.style.display = (elt.style.display == "table-row" ? "none" : "table-row");
}