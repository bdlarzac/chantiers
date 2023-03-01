<?php
/******************************************************************************
    Calcule la liste des essences présentes dans les ugs
    
    @license    GPL
    @history    2023-02-28 20:26:25+01:00, Thierry Graff : Creation
********************************************************************************/

require_once 'lib-php/csvAssociative.php';

$csv = csvAssociative::compute('../data/ug.csv', ',');

compute_essence($csv);

/** 
**/
function compute_essence(&$csv){
    $res = [];
    foreach($csv as $row){
        $strEssences = trim($row['Essence']);
        $essences = explode('+', $strEssences);
        foreach($essences as $ess){
            $ess = trim($ess);
            if(!isset($res[$ess])){
                $res[$ess] = 'bidon';
            }
        }
    }
    ksort($res);
    foreach(array_keys($res) as $ess){
        echo "$ess\n";
    }
}
