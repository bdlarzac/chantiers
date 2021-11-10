/******************************************************************************
    Code relatif aux chantiers plaquettes
    Mis ici car commun à plusieurs pages

    @license    GPL
    @history    2020-01-27 17:57:44+01:00, Thierry Graff : Creation
********************************************************************************/

function deleteChantierPlaquette(idChantier, nomChantier){
    const msg = "ATTENTION, en cliquant sur OK,\n"
            + "le chantier plaquette \"" + nomChantier + "\" sera définitivement supprimé.\n"
            + "\nToutes les opérations associées à ce chantier seront aussi supprimées :"
            + "\n- Abattage,"
            + "\n- Débardage,"
            + "\n- Broyage,"
            + "\n- Déchiquetage,"
            + "\n- Transport plateforme,"
            + "\n- Rangement."
            + "\nLes tas associés au chantier seront aussi supprimés."
            + "\nCette opération est définitive ET NE PEUT PAS ETRE ANNULEE.";
    const r = confirm(msg);
    if (r == true) {
        window.location = "/chantier/plaquette/delete/" + idChantier;
    }
}

function quantiteBoisSec(qteVert, pourcentPerte){
    return qteVert * (1- pourcentPerte / 100);
}