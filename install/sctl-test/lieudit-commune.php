<?php
/******************************************************************************
    Teste les associations lieux-dits / communes
    Utilise les csv extraits de la base SCTL avec mdb-tools
    
    @license    GPL
    @history    2019-11-03 10:45:15+01:00, Thierry Graff : Creation
********************************************************************************/

require_once 'lib-php/csvAssociative.php';

// Tables de la base SCTL où on trouve lieux-dits et communes :
$tables = [
    // 'PSG' => ['', ''], // IdLieuDit, nom commune
    'Subdivision' => ['IdCommune', 'IdLieuDit'],
    'SubdivCadastre' => ['IdCommune', 'IdLieuDit'],
    // 'ParcelleSup' => ['', ''], id commune, nom lieu dit
    'Parcelle' => ['IdCommune', 'IdLieuDit'],
    // ForPrintSctl (pas id)
];

$config = yaml_parse(file_get_contents('../../config.yml'));

//$version = '2018';
//$version = '2020-02-27';
$version = '2020-12-16';

$csvDir = $config['dev']['sctl-data'] . "/csv-$version";

//test_1n($csvDir, $tables);
//test_pareil($csvDir, $tables);
//test_coherence($csvDir, $tables);
//liste_multiples($csvDir, $tables);
test_zero($csvDir, $tables);

// ******************************************************
/**
    Teste s'il existe des lieux-dits sans communes
    Utilise SubdivCadastre
    Résultat : 
    2018 : tous les lieux-dits sont associés à une commune.
    2020-02-27
    2020-12-16 :
        COMBEGRAND
        LE PERTUS
        LES AUMIERES
        PUECH BLACOUS
        PUECH ROUCOUS ET PEY
        REDOULES
        ROUTAOUS
        SERRE DE LA BAUME
        SERRE DE LAS SOUTS
        => 9 lieux-dits sans commune
**/
function test_zero($csvDir, $tables){
    echo "Teste lieux-dits sans communes\n";
    $table = 'SubdivCadastre';
    $res = [];
    $nameIdC = $tables[$table][0];
    $nameIdLD = $tables[$table][1];
    $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
    foreach($csv as $row){
        $idC = $row[$nameIdC];
        $idLD = $row[$nameIdLD];
        if(!isset($res[$idLD])){
            $res[$idLD] = [];
        }
        $res[$idLD][] = $idC;
    }
    
    // remove doublons
    foreach($res as $idLD => $idsC){
        $idsC = array_values(array_unique($idsC)); // array_values to reindex
        $res[$idLD] = $idsC;
    }
    
    $communes = []; // assoc id - nom
    $csv = csvAssociative::compute($csvDir . '/Commune.csv');
    foreach($csv as $row){
        $communes[$row['COMMUNE']] = $row['NOM'];
    }
    
    $lieuxdits = []; // assoc id - nom
    $csv = csvAssociative::compute($csvDir . '/LieuDit.csv');
    foreach($csv as $row){
        $lieuxdits[$row['IdLieuDit']] = $row['Libelle'];
    }
    
    // test 0 associations
    $N = 0;
    foreach($lieuxdits as $idLD => $nomLD){
        if(!isset($res[$idLD])){
            echo "$nomLD\n";
            $N++;
        }
    }
    echo "=> $N lieux-dits sans commune\n";
}

// ******************************************************
/**
    Liste les lieux-dits associés à plusieurs communes
    Utilise SubdivCadastre
**/
function liste_multiples($csvDir, $tables){
    $table = 'SubdivCadastre';
    $res = [];
    $nameIdC = $tables[$table][0];
    $nameIdLD = $tables[$table][1];
    $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
    foreach($csv as $row){
        $idC = $row[$nameIdC];
        $idLD = $row[$nameIdLD];
        if(!isset($res[$idLD])){
            $res[$idLD] = [];
        }
        $res[$idLD][] = $idC;
    }
    
    // remove doublons
    foreach($res as $idLD => $idsC){
        $idsC = array_values(array_unique($idsC)); // array_values to reindex
        $res[$idLD] = $idsC;
    }
    
    $communes = []; // assoc id - nom
    $csv = csvAssociative::compute($csvDir . '/Commune.csv');
    foreach($csv as $row){
        $communes[$row['COMMUNE']] = $row['NOM'];
    }
    
    $lieuxdits = []; // assoc id - nom
    $csv = csvAssociative::compute($csvDir . '/LieuDit.csv');
    foreach($csv as $row){
        $lieuxdits[$row['IdLieuDit']] = $row['Libelle'];
    }
    
    $N = 0;
    ksort($res);
    foreach($res as $idLD => $idsC){
        if(count($idsC) == 1){
            continue;
        }
        $N++;
        echo "\nLieu-dit $idLD " . $lieuxdits[$idLD] . "\n";
        foreach($idsC as $idC){
            echo "    $idC " . $communes[$idC] . "\n";
        }
    }
    echo "\n$N lieux-dits concernés\n";
}

// ******************************************************
/**
    Teste si les 3 tables ont les mêmes communes et les mêmes lieudit
    Résultat :
        communes : ok, les 3 tables sont pareil
        lieudit : lieudit 393 (TRAVERS DE CLAPADE) pas dans Parcelle
                  les 2 autres tables contiennent les mêmes lieudit
        => Pour fabriquer la table de liens, utiliser SubdivCadastre ou Subdivision, mais pas Parcelle.
**/
function test_pareil($csvDir, $tables){
    
    // teste id commune
    $idsC = [];
    foreach($tables as $table => $details){
        $idsC[$table] = [];
        $nameIdC = $details[0];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $idC = $row[$nameIdC];
            $idsC[$table][] = $row[$nameIdC];
        }
    }
    foreach($idsC as $table => $elements){
        $idsC[$table] = array_values(array_unique($elements));
    }
    
    echo "Compare id communes :\n";
    [$t1, $t2, $t3] = array_keys($tables);
    //
    $test = array_diff($idsC[$t1], $idsC[$t2]);
    echo "  $t1 - $t2 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t1 et pas dans $t2 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsC[$t2], $idsC[$t1]);
    echo "  $t2 - $t1 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t2 et pas dans $t1 : " . implode(' ', $test) . "\n";
    }
    //
    $test = array_diff($idsC[$t1], $idsC[$t3]);
    echo "  $t1 - $t3 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t1 et pas dans $t3 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsC[$t3], $idsC[$t1]);
    echo "  $t3 - $t1 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t3 et pas dans $t1 : " . implode(' ', $test) . "\n";
    }
    //
    $test = array_diff($idsC[$t2], $idsC[$t3]);
    echo "  $t2 - $t3 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t2 et pas dans $t3 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsC[$t3], $idsC[$t2]);
    echo "  $t3 - $t2 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t3 et pas dans $t2 : " . implode(' ', $test) . "\n";
    }
    
    // teste id lieudit
    $idsLD = [];
    foreach($tables as $table => $details){
        $idsLD[$table] = [];
        $nameIdLD = $details[1];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $idLD = $row[$nameIdLD];
            $idsLD[$table][] = $row[$nameIdLD];
        }
    }
    foreach($idsLD as $table => $elements){
        $idsLD[$table] = array_values(array_unique($elements));
    }
    
    echo "Compare id lieudit :\n";
    [$t1, $t2, $t3] = array_keys($tables);
    //
    $test = array_diff($idsLD[$t1], $idsLD[$t2]);
    echo "  $t1 - $t2 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t1 et pas dans $t2 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsLD[$t2], $idsLD[$t1]);
    echo "  $t2 - $t1 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t2 et pas dans $t1 : " . implode(' ', $test) . "\n";
    }
    //
    $test = array_diff($idsLD[$t1], $idsLD[$t3]);
    echo "  $t1 - $t3 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t1 et pas dans $t3 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsLD[$t3], $idsLD[$t1]);
    echo "  $t3 - $t1 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t3 et pas dans $t1 : " . implode(' ', $test) . "\n";
    }
    //
    $test = array_diff($idsLD[$t2], $idsLD[$t3]);
    echo "  $t2 - $t3 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t2 et pas dans $t3 : " . implode(' ', $test) . "\n";
    }
    $test = array_diff($idsLD[$t3], $idsLD[$t2]);
    echo "  $t3 - $t2 : ";
    if(count($test) == 0){
        echo "OK\n";
    }
    else{
        echo "\n    Dans $t3 et pas dans $t2 : " . implode(' ', $test) . "\n";
    }
    
    echo "\n";
}


// ******************************************************
/**
    Teste si Subdivision et SubdivCadastre ont les mêmes associations commune - lieudit
    Résultat : pas de différence, tous les couples (commune, lieudit) sont pareils dans les 2 tables.
**/
function test_coherence($csvDir, $tables){
    $res = [];
    foreach($tables as $table => $details){
        $res[$table] = [];
        $nameIdC = $details[0];
        $nameIdLD = $details[1];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $idC = $row[$nameIdC];
            $idLD = $row[$nameIdLD];
            if(!isset($res[$table][$idLD])){
                $res[$table][$idLD] = [];
            }
            $res[$table][$idLD][] = $idC;
        }
    }
    
    // remove doublons
    foreach($res as $table => $infos){
        foreach($infos as $idLD => $idsC){
            $idsC = array_values(array_unique($idsC)); // array_values to reindex
            $res[$table][$idLD] = $idsC;
        }
    }
    
    // on ne compare que pour Subdivision et SubdivCadastre
    // puisque test_pareil() a montré qu'il faut éviter d'utiliser Parcelle
    [$t1, $t2, $t3] = array_keys($tables);
    
    echo "Compare associations $t1 - $t2\n";
    $ok = true;
    foreach($res[$t1] as $idLD => $idsC){
        $test = array_diff($idsC, $res[$t2][$idLD]);
        if(count($test) != 0){
            echo "\n    Lieu-dit $idLD : Dans $t1 et pas dans $t2 : " . implode(' ', $test) . "\n";
            $ok = false;
        }
    }
    if($ok){
        echo "    OK, pas de différence\n";
    }
    
    
    echo "Compare associations $t2 - $t1\n";
    $ok = true;
    foreach($res[$t2] as $idLD => $idsC){
        $test = array_diff($idsC, $res[$t1][$idLD]);
        if(count($test) != 0){
            echo "\n    Lieu-dit $idLD : Dans $t2 et pas dans $t1 : " . implode(' ', $test) . "\n";
            $ok = false;
        }
    }
    if($ok){
        echo "    OK, pas de différence\n";
    }
    
    
}

// *********************************************************
/** 
    Teste s'il y a un lien 1-n entre lieux-dits et communes.
    Résultat : non, lien n-n pour 24 lieux-dits
**/
function test_1n($csvDir, $tables){
    
    $res = [];
    foreach($tables as $table => $details){
        $res[$table] = [];
        $nameIdC = $details[0];
        $nameIdLD = $details[1];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $idC = $row[$nameIdC];
            $idLD = $row[$nameIdLD];
            if(!isset($res[$table][$idLD])){
                $res[$table][$idLD] = [];
            }
            $res[$table][$idLD][] = $idC;
        }
    }
    
    $multiples = [];
    foreach($res as $table => $infos){
        $multiples[$table] = [];
        foreach($infos as $idLD => $idsC){
            $idsC = array_values(array_unique($idsC)); // array_values to reindex
            if(count($idsC) > 1){
                $multiples[$table][$idLD] = $idsC;
            }
        }
    }
    
    $lisibles = []; // keys = idLD
    foreach($multiples as $table => $multiple){
        foreach($multiple as $idLD => $idsC){
            $lisibles[$idLD][$table] = $idsC;
        }
    }
    
    foreach($lisibles as $idLD => $details){
        echo "Lieu dit $idLD\n";
        foreach($details as $table => $idsC){
            echo "    " . str_pad($table, 15) . ' : ' . implode(' ', $idsC) . "\n";
        }
    }
    echo count($lisibles) . " cas multiples\n";
}
