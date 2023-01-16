<?php
/******************************************************************************
    Détecte les parcelles avec le même code à 6 caractères.
    Pour fixer un bug lors de la création de la base.
    Le code 6 est unique au sein d'une commune.
    Issue #11
    
    Résultat sur la prod du 2023-01-11 :
    Problème potentiel uniquement pour chantiers chauffage fermier
    Seul le chantier 5 a un problème avec la parcelle 0N0321 (sur Millau et La Couvertoirade)
    
    @license    GPL
    @history    2023-01-11 17:33:33+01:00, Thierry Graff : Creation
********************************************************************************/

require_once 'lib-php/csvAssociative.php';

$csv = csvAssociative::compute('../sctl-data/csv-2023-01-11/Parcelle.csv', ';');

teste_parcelles($csv);

function teste_parcelles(&$csv){
    //
    // 1 - Parcelles en double
    //
    $all = []; // assoc code parcelle => nb occurence
    foreach($csv as $line){
        $codeP = $line['PARCELLE'];
        if(!isset($all[$codeP])){
            $all[$codeP] = 1;
        } else {
            $all[$codeP]++;
        }
    }
    echo "Parcelles en double :\n";
    $N = 0;
    foreach($all as $code => $n){
        if($n > 1){
            echo "$code : $n\n";
            $N++;
        }
    }
    ksort($all);
    echo "----------\n";
    echo "$N doubles\n";
    echo "----------\n";
    //
    // 2 - Parcelles problématiques
    //
    // sur la base de prod le 2023-01-11
    // select code from parcelle where id in(select id_parcelle from chaufer_parcelle);
    $prods = [
        '0G0179',
        '0N0321',
        '0D0042',
        '0K0109',
        '0N0042',
        '0N0044',
        '0E0146',
        '0G0111',
        '0G0178',
        '0N0053',
    ];
    foreach($prods as $prod){
        if($all[$prod] > 1){
            echo "PROBLEM sur la prod : $prod\n";
        }
    }
    
    
}
