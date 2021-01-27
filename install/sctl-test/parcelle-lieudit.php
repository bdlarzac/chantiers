<?php
/******************************************************************************
    Teste les associations parcelles / lieux-dits
    Utilise les csv extraits de la base de terres avec mdb-tools

    @history    2019-11-06 06:41:27+01:00, Thierry Graff : Creation
********************************************************************************/

/* 
Tables contenant des associations :

Pas utilisables (car ne contient pas ids) :
ForPrintSctl
    PARCELLE (code string)
    Lieudit (nom)
PSG
	IdParcelle (double) - à priori transformable en int
	Parcelle (code string)
	LieuDit (nom) 
ParcelleSup
    PARCELLE (code string)
    LIEUDIT (nom)
    IdLieuDit ; parfois pas renseigné

Utilisables : 
Parcelle
    PARCELLE (code string)
    IdParcelle (int)
    IdLieuDit
SubdivCadastre
	IdParcelle (int)
	IdLieuDit
Subdivision
	IdParcelle (int)
	IdLieuDit
	
Remarques :
Une parcelle est caractérisée par
    - code (ex 0B0206)
    - id (int)
Cette association est uniquement présente dans les tables Parcelle et PSG

Conclusion après exec des tests :
On utilise Parcelle pour construire
- une table parcelle
    id int pk
    code string
- une table parcelle_lieudit
    id_parcelle int fk
    id_lieudit int fk
Pour la prod, le pb de \n qui oblige à faire Parcelle ne se posera pas
si on part de la base access à la place du csv    
*/

require_once 'lib-php/csvAssociative.php';

$config = yaml_parse(file_get_contents('../../config.yml'));

//$version = '2018';
//$version = '2020-02-27';
$version = '2020-12-16';

$csvDir = $config['dev']['sctl-data'] . "/csv-$version";

$tables = ['Parcelle', 'SubdivCadastre', 'Subdivision'];

test_ids($csvDir, $tables);
//test_idCodeParcelle($csvDir);
//test_assoc($csvDir, $tables);
//test_surface($csvDir);


// ******************************************************
/**
    Affiche surface min et max des parcelles dans Parcelle
    Résultats :
        min = 7.0000000000000000e+00
        max = 1.5183800000000000e+06
**/
function test_surface($csvDir){
    $csv = csvAssociative::compute($csvDir . '/Parcelle.csv');
    $min = PHP_INT_MAX;
    $max = PHP_INT_MIN;
    foreach($csv as $row){
        $s = $row['SURFACE'];
        if($s > $max){
            $max = $s;
        }
        if($s < $min){
            $min = $s;
        }
    }
    echo "min = $min\n";
    echo "max = $max\n";
}

// ******************************************************
/**
    Sachant que test_ids() a montré que les 3 tables ont les mêmes ids parcelle,
    teste les associations entre id parcelle et id lieudit
    Résultats :
        - Aucune différence entre les 3 tables
**/
function test_assoc($csvDir, $tables){
    $data = [];
    foreach($tables as $table){
        $data[$table] = [];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $idP = $row['IdParcelle'];
            if(!isset($data[$table][$idP])){
                $data[$table][$idP] = [];
            }
            $data[$table][$idP][] = $row['IdLieuDit'];
        }
    }
    foreach($tables as $table){
        foreach($data[$table] as $idP => $idsLD){
            $data[$table][$idP] = array_unique($idsLD);
        }
    }
    
    $t1 = $tables[0];
    $t2 = $tables[1];
    foreach($data[$t1] as $idP => $idsLD){
        $idsLD2 = $data[$t2][$idP];
        $diff1 = array_diff($idsLD, $idsLD2);
        if(count($diff1) != 0){
            echo "problem\n";
        }
        $diff2 = array_diff($idsLD2, $idsLD);
        if(count($diff2) != 0){
            echo "problem\n";
        }
    }
    $t1 = $tables[0];
    $t2 = $tables[2];
    foreach($data[$t1] as $idP => $idsLD){
        $idsLD2 = $data[$t2][$idP];
        $diff1 = array_diff($idsLD, $idsLD2);
        if(count($diff1) != 0){
            echo "problem\n";
        }
        $diff2 = array_diff($idsLD2, $idsLD);
        if(count($diff2) != 0){
            echo "problem\n";
        }
    }
    $t1 = $tables[1];
    $t2 = $tables[2];
    foreach($data[$t1] as $idP => $idsLD){
        $idsLD2 = $data[$t2][$idP];
        $diff1 = array_diff($idsLD, $idsLD2);
        if(count($diff1) != 0){
            echo "problem\n";
        }
        $diff2 = array_diff($idsLD2, $idsLD);
        if(count($diff2) != 0){
            echo "problem\n";
        }
    }
}


// ******************************************************
/**
    Teste les associations id parcelle - code parcelle dans Parcelle et PSG
    Résultats :
    - Chaque id est associé à un seul code
    - un code peut être associé à plusieurs ids (202 codes dans ce cas)
    => id = clé primaire ; code = caractéristique
    
    - Les associations ids <=> codes sont exactement les mêmes dans les 2 tables
**/
function test_idCodeParcelle($csvDir){
    
    $tables = [
        'Parcelle' => ['IdParcelle', 'PARCELLE'],
        'PSG' => ['IdParcelle', 'Parcelle'],
    ];
    
    $idCodes = []; // id => code
    $codeIds = []; // code => id
    foreach($tables as $table => $details){
        echo "=== $table ===\n";
        $idName = $details[0];
        $codeName = $details[1];
    
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        
        $idCodes[$table] = []; // id => code
        $codeIds[$table] = []; // code => id
        foreach($csv as $row){
            $id = (int)$row[$idName];
            $code = trim($row[$codeName]);
            if(!isset($idCodes[$table][$id])){
                $idCodes[$table][$id] = [];
            }
            $idCodes[$table][$id][] = $code;
            
            if(!isset($codeIds[$table][$code])){
                $codeIds[$table][$code] = [];
            }
            $codeIds[$table][$code][] = $id;
        }
        
        foreach(array_keys($idCodes[$table]) as $id){
            $idCodes[$table][$id] = array_unique($idCodes[$table][$id]);
        }
        foreach(array_keys($codeIds[$table]) as $code){
            $codeIds[$table][$code] = array_unique($codeIds[$table][$code]);
        }
        
        // test assoc id <=> code
        echo "Test id => code\n";
        $N = 0;
        foreach($idCodes[$table] as $id => $codes){
            if(count($codes) != 1){
                echo $id . implode(' ', $codes) . "\n";
                $N++;
            }
        }
        echo "Assoc multiples id => codes : $N\n";
        
        echo "Test code => id\n";
        $N = 0;
        foreach($codeIds[$table] as $code => $ids){
            if(count($ids) != 1){
                echo $code . ' : ' . implode(' ', $ids) . "\n";
                $N++;
            }
        }
        echo "Assoc multiples code => ids : $N\n";
    }             
    
    // compare assoc dans les deux tables
    foreach($idCodes['Parcelle'] as $id => $codes){
        if(!isset($idCodes['PSG'][$id])){
            echo "Manque id $id dans PSG\n";
        }
        $codes2 = $idCodes['PSG'][$id];
        $diff1 = array_diff($codes, $codes2);
        if(count($diff1) != 0){
            echo "\n"; print_r($diff1); echo "\n";
        }
        $diff2 = array_diff($codes2, $codes);
        if(count($diff2) != 0){
            echo "\n"; print_r($diff2); echo "\n";
        }
    }
    foreach($idCodes['PSG'] as $id => $codes){
        if(!isset($idCodes['Parcelle'][$id])){
            echo "Manque id $id dans Parcelle\n";
        }
        $codes2 = $idCodes['Parcelle'][$id];
        $diff1 = array_diff($codes, $codes2);
        if(count($diff1) != 0){
            echo "\n"; print_r($diff1); echo "\n";
        }
        $diff2 = array_diff($codes2, $codes);
        if(count($diff2) != 0){
            echo "\n"; print_r($diff2); echo "\n";
        }
    }
}

// ******************************************************
/**
    Teste si les 3 tables utilisables ont les mêmes ids parcelle
    Résultat :
        les 3 tables ont exactement les mêmes ids
        2361 ids distincts
**/
function test_ids($csvDir, $tables){
    $data = [];
    foreach($tables as $table){
        $data[$table] = [];
        $csv = csvAssociative::compute($csvDir . '/' . $table . '.csv');
        foreach($csv as $row){
            $data[$table][] = $row['IdParcelle'];
        }
    }
    foreach($tables as $table){
        $data[$table] = array_unique($data[$table]);
        echo $table . ' : ' . count($data[$table]) . "\n";
    }
    
    $diff1 = array_diff($data[$tables[0]], $data[$tables[1]]);
    $str1 = "Dans {$tables[0]} et pas dans {$tables[1]}";
    $diff1a = array_diff($data[$tables[1]], $data[$tables[0]]);
    $str1a = "Dans {$tables[1]} et pas dans {$tables[0]}";
    //
    $diff2 = array_diff($data[$tables[0]], $data[$tables[2]]);
    $str2 = "Dans {$tables[0]} et pas dans {$tables[2]}";
    $diff2a = array_diff($data[$tables[2]], $data[$tables[0]]);
    $str2a = "Dans {$tables[2]} et pas dans {$tables[0]}";
    //
    $diff3 = array_diff($data[$tables[1]], $data[$tables[2]]);
    $str3 = "Dans {$tables[1]} et pas dans {$tables[2]}";
    $diff3a = array_diff($data[$tables[2]], $data[$tables[1]]);
    $str3a = "Dans {$tables[2]} et pas dans {$tables[1]}";
    echo "$str1 : " . count($diff1) . "\n";
    echo "$str1a : " . count($diff1a) . "\n";
    echo "\n";
    echo "$str2 : " . count($diff2) . "\n";
    echo "$str2a : " . count($diff2a) . "\n";
    echo "\n";
    echo "$str3 : " . count($diff3) . "\n";
    echo "$str3a : " . count($diff3a) . "\n";
}
