<?php

try{
    //$dbhAccess = new PDO("odbc:Driver={Microsoft Access Driver (*.mdb, *.accdb)};Dbq=/home/thierry/dev/jobs/bdl/2.build/dbs/2-Sctl-Gfa.mdb;Uid=Admin");
    $dbhAccess = new PDO("odbc:Driver=MDBTools;DBQ=/home/thierry/dev/jobs/bdl/2.build/dbs/2-Sctl-Gfa.mdb;");
}
catch(PDOException $e){
    echo $e->getMessage() . "\n";
    $e->printStackTrace();
    exit();
}

