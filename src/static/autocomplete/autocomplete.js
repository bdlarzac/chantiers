/******************************************************************************
        Adaptation from https://www.w3schools.com/howto/howto_js_autocomplete.asp

        @license        GPL
        @history        2020-01-17 16:35:45+01:00, Thierry Graff : Creation
********************************************************************************/


/**
    Generic function which adds html and css to a input type=text field.
    @param inputField     The input text field element on which autocomplete is added
    @param dataProvider   Function returning a regular array
                          Provides the data to use to fill autocomplete values
    Note : le test commenté permet d'avoir "LE PINEL" lorsqu'on tape "PI"
        Le test est utile dans l'exemple de w3school car arr contient tout.
        Inutile ici car arr contient seulement des entrées correspondant à la saisie.
**/
function autocomplete(inputField, dataProvider) {
    var currentFocus;
    inputField.addEventListener(
        "input",
        function(e) {
            var a, b, i, val = this.value;
            closeAllLists();
            // 2 could be passed in param (minimal length before triggering autocomplete)
            if (!val || val.length < 2) {
                return false;
            }
            currentFocus = -1;
            a = document.createElement("DIV");
            a.setAttribute("id", this.id + "autocomplete-list");
            a.setAttribute("class", "autocomplete-items");
            this.parentNode.appendChild(a);
            // here uses function passed in parameter to retrieve values (in general with ajax)
            arr = dataProvider(inputField);
            for (i = 0; i < arr.length; i++) {
                // voir note
                //if (arr[i].substr(0, val.length).toUpperCase() == val.toUpperCase()) {
                    b = document.createElement("DIV");
                    // version originale, qui met en gras les 2 premiers caractères
                    //b.innerHTML = "<strong>" + arr[i].substr(0, val.length) + "</strong>";
                    //b.innerHTML += arr[i].substr(val.length);
                    // version modifiée, qui met en gras le texte saisi ; pas forcément les 2 premiers caractères
                    // Attention, bug con si on utilise strong à la place de b :
                    // si val="st" ou "rt" ou "ro" etc. le 2e replace va remplacer à l'intérieur de <strong> 
                    b.innerHTML = arr[i].replace(val.toUpperCase(), "<b>" + val.toUpperCase() + "</b>")
                           .replace(val.toLowerCase(), "<b>" + val.toLowerCase() + "</b>");
                    b.innerHTML += "<input type='hidden' value='" + arr[i] + "'>";
                            b.addEventListener("click", function(e) {
                            inputField.value = this.getElementsByTagName("input")[0].value;
                            closeAllLists();
                    });
                    a.appendChild(b);
                //}
            }
        }
    );
    
    inputField.addEventListener(
        "keydown",
        function(e) {
            var x = document.getElementById(this.id + "autocomplete-list");
            if (x) x = x.getElementsByTagName("div");
            if (e.keyCode == 40) { // down
                currentFocus++;
                addActive(x);
            } else if (e.keyCode == 38) { // up
                currentFocus--;
                addActive(x);
            } else if (e.keyCode == 27) { // escape
                closeAllLists();
            } else if (e.keyCode == 13) { // enter
                e.preventDefault();
                if (currentFocus > -1) {
                    if (x) x[currentFocus].click();
                }
            }
        }
    );
    
    function addActive(x) {
        if (!x) return false;
        removeActive(x);
        if (currentFocus >= x.length) currentFocus = 0;
        if (currentFocus < 0) currentFocus = (x.length - 1);
        x[currentFocus].classList.add("autocomplete-active");
    }
    
    function removeActive(x) {
        for (var i = 0; i < x.length; i++) {
            x[i].classList.remove("autocomplete-active");
        }
    }
    
    function closeAllLists(elmnt) {
        var x = document.getElementsByClassName("autocomplete-items");
        for (var i = 0; i < x.length; i++) {
            if (elmnt != x[i] && elmnt != inputField) {
                x[i].parentNode.removeChild(x[i]);
            }
        }
    }

    document.addEventListener(
        "click",
        function (e) {
            closeAllLists(e.target);
        }
    );
}
