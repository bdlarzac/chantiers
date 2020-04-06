/******************************************************************************
    Code relatif aux chantiers plaquettes
    Mis ici car commun à plusieurs pages

    @license    GPL
    @history    2020-01-27 17:57:44+01:00, Thierry Graff : Creation
********************************************************************************/


// *****************************************
function deleteChantierPlaquette(idChantier, nomChantier){
    var msg = "ATTENTION, en cliquant sur OK,\n"
            + "le chantier plaquette \"" + nomChantier + "\" sera définitivement supprimé.\n"
            + "\nToutes les opérations associées à ce chantier seront aussi supprimées :\n"
            + "abattage, débardage, broyage, déchiquetage, transport plateforme, rangement\n";
    var r = confirm(msg);
    if (r == true) {
        window.location = "/chantier/plaquette/delete/" + idChantier;
    }
}
