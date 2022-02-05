/******************************************************************************
    Formats a number with grouped thousands.
    Useless in theory, as toLocaleString() is supposed to do the job
    - but didn't succeed to toLocaleString() correctly.
    Ex: 
        formatNb(2000) = "2 000"
        formatNb(2000.12) = "2 000.12"
        formatNb(1232000.12) = "1 232 000.12"
        formatNb(1232000) = "1 232 000"
        formatNb(20.12) = "20.12"
        formatNb(20) = "20"
    @param  x Number to format
    @param  sep Character used to separate thousands. 
    @param  sepIntDec Character separating integer and decimal parts of the number.
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
    
    @history    2021-12-06 18:52:34+01:00, Thierry Graff : Creation
********************************************************************************/
function formatNb(x, sep=' ', sepIntDec='.'){
    [intPart, decPart] = x.toString().split(sepIntDec);
    let res = intPart[0];
    for(let i=1; i < intPart.length ; i++){
        if((intPart.length - i)%3 == 0){
            res = res + sep;
        }
        res += intPart[i];
    }
    if(decPart != undefined){
        res += sepIntDec + decPart;
    }
    return res;
}
