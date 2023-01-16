<?php
/******************************************************************************
    Teste les associations parcelles / ug
    Teste les caractéristiques de ug
    Utilise un export de gis.
    
    Résultats : on a une association n-n
    
    @license    GPL
    @history    2019-10-26 23:53:22+02:00, Thierry Graff : Creation
********************************************************************************/

require_once 'csvAssociative.php';

$csv = csvAssociative::compute('/home/thierry/dev/jobs/bdl/2.build/gis/mapinfo-ug-test.csv', ',');

//ug_parcelles($csv);
//parcelles_ug($csv);
//carac_ug($csv);
//teste_coupe($csv);
teste_essence($csv);

// ***********************************************
/** 
    Regarde si, pour un code ug donné, il y a plusieurs valeurs différentes
    des essences (types de peuplement)
    Résultats : unique
**/
function teste_essence(&$csv){
    $res = [];
    foreach($csv as $row){
        $code = $row['PG'];
        if($code == ''){
            continue;
        }
        $essence = trim($row['Essence']);
        if(!isset($res[$code])){
            $res[$code] = [];
        }
        $res[$code][] = $essence;
    }
    foreach($res as $code => $values){
        $res[$code] = array_unique($values);
        if(count($res[$code]) != 1){
            echo "$code : '" . implode("' '", $res[$code]) . "'\n";
        }
    }
}


// ***********************************************
/** 
    Regarde si, pour un code ug donné, il y a plusieurs valeurs différentes de
    (type_coupe, annee_intervention, previsionnel_coupe)
    Résultats :
        Plusieurs valeurs différentes pour
            IX-21 : 'ESP2 2020+ESP3' 'ESP2 0+ESP3'
            V-5 : 'E1 2019+E2' 'E1 2021+E2'
            XVI-12 : 'Degagement 2015+Dpressage?' 'Degagement 2016+Dpressage?'
            XVIII-9 : 'ESP3 2020+' 'ESP3 0+'
            XVIII-5 : 'ESP3 0+' 'ESP3 2020+'
            XVII-1 : 'ESP3 2020+ESP3' 'ESP3 0+ESP3'
            XIV-9 : ' 0+E2' 'E1 2016+E2'
            IX-11 : 'CR 2020+0' 'ESP2 2014+0'
        Même dans ces cas, previsionnel_coupe est unique
        Mais (type_coupe, annee_intervention) n'est pas unique
**/
function teste_coupe(&$csv){
    $res = [];
    foreach($csv as $row){
        $code = $row['PG'];
        if($code == ''){
            continue;
        }
        $type_coupe = trim($row['Coupe']) . ' ' . trim($row['Annee_intervention']);
        $previsionnel_coupe = trim($row['PSG_suivant']);
        if(!isset($res[$code])){
            $res[$code] = [];
        }
        $res[$code][] = "$type_coupe+$previsionnel_coupe";
    }
    foreach($res as $code => $values){
        $res[$code] = array_unique($values);
        if(count($res[$code]) != 1){
            echo "$code : '" . implode("' '", $res[$code]) . "'\n";
        }
    }
}


// ***********************************************
/** 
    Résultats :
    Nb codes ug = 660
    Max len(code ug) = 8
**/
function carac_ug(&$csv){
    $res = [];
    foreach($csv as $row){
        $res[] = $row['PG'];
    }
    $res = array_values(array_unique($res));
    $max = 0;
    foreach($res as $code){
        $len = strlen($code);
        if($len > $max){
            $max = $len;
        }
    }
    echo "Nb codes ug = " . count($res) . "\n";
    echo "Max len(code ug) = $max\n";
}


// ***********************************************
/** 
    Montre ug => parcelles
**/
function ug_parcelles(&$csv){
    $res = [];
    foreach($csv as $row){
        $id = $row['PG'];
        if(!isset($res[$id])){
            $res[$id] = [];
        }
        $res[$id][] = substr($row['ID_PARCELLE_11'], 5);
    }
    foreach($res as $k => $v){
        $res[$k] = array_unique($v);
    }
    echo "\n"; print_r($res); echo "\n";
}


// ***********************************************
/** 
    Montre parcelles => ug
**/
function parcelles_ug(&$csv){
    $res = [];
    foreach($csv as $row){
        $id = substr($row['ID_PARCELLE_11'], 5);
        if(!isset($res[$id])){
            $res[$id] = [];
        }
        $res[$id][] = $row['PG'];
    }
    foreach($res as $k => $v){
        $res[$k] = array_unique($v);
    }
    echo "\n"; print_r($res); echo "\n";
}