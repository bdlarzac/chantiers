/*
Fixe certains parcelle manquante et associations avec UGs

Ce fix est effectué sur les ugs du PSG1

Copie de l'analyse du problème sur https://github.com/bdlarzac/chantiers/issues/27 :

-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
Lors de la création d'un nouveau chantier plaquettes, message reçu

ERREUR - Transmettez ce message l'administrateur du site :
Mauvais format de retour de /ajax/get/parcelles-from-ids-ugs/282,513

Retour ajax envoyé par afficheParcelles() dans src/view/common/liens-parcelles.html

-------------------------------------------------------------------------------
-------------------------------------------------------------------------------

Les UGs existent bien :

bdlchantiers=> select id,code from ug where id in(282,513);
 id  |  code  
-----+--------
 282 | VII-8
 513 | VII-22   

Mais elles ne sont asssociées à aucune parcelle :

bdlchantiers=> select * from parcelle where id in(select id_parcelle from parcelle_ug where id_ug in(282,513)) order by code;
 id | id_proprietaire | code | surface | id_commune 
----+-----------------+------+---------+------------
(0 rows)

-------------------------------------------------------------------------------
-------------------------------------------------------------------------------
Les UGs existent bien :

bdlchantiers=> select id,code from ug where id in(282,513);
 id  |  code  
-----+--------
 282 | VII-8
 513 | VII-22   

Mais elles ne sont asssociées à aucune parcelle :

bdlchantiers=> select * from parcelle where id in(select id_parcelle from parcelle_ug where id_ug in(282,513)) order by code;
 id | id_proprietaire | code | surface | id_commune 
----+-----------------+------+---------+------------
(0 rows)

Ce problème concerne 5 UGs :

bdlchantiers=> select id,code from ug where id not in(select id_ug from parcelle_ug); 
 id  |  code  
-----+--------
  33 | VII-27
 282 | VII-8
 624 | VII-16
 513 | VII-22
 563 | I-53

Heureusement, ces parcelles ne sont associées à aucun chantier :

bdlchantiers=> select * from chantier_ug where id_ug in (select id from ug where id not in(select id_ug from parcelle_ug)); 
 type_chantier | id_chantier | id_ug 
---------------+-------------+-------
(0 rows)

En vérifiant dans le fichier manage/data/ug.csv ayant servi à la création de la base, on peut voir que ces UGs sont pourtant liées à des parcelles, à part l'UG I-53 :

id_ug | code_ug     | parcelles correspondantes
------+-------------+--------------------------
563   | I-53        0
282   | VII-8       121450M0168
624   | VII-16      121450M0168
513   | VII-22      121450M0168
 33   | VII-27      121450M0168

Donc bug datant de la création de la base, passé inaperçu jusqu'à maintenant car ces UGs ne sont impliquées dans aucun chantier.    
Tous les liens manquant concernent la même parcelle : 121450M0168
Parcelle présente dans la base SCTL de 2018 mais pas dans les versions suivantes (à partir de 2020-02-27)

-------------------------------------------------------------------------------
Solution retenue
-------------------------------------------------------------------------------
Ajout de la parcelle et de ses liens - pas réussi à savoir pourquoi cette parcelle était absente

@copyright  BDL, Bois du Larzac
@license    GPL
@history    2024-04-17 08:35:38+02:00, Thierry Graff : Creation
*/
package main

import (
	"bdl.local/bdl/ctxt"
	"fmt"
)

func Migrate_2024_04_17_ugs_orphelines__27(ctx *ctxt.Context) {
	db := ctx.DB
	var query string
	//
	// ajout parcelle 0M0168
	//
    // table parcelle : id n'est pas serial (prend la valeur dans la base sctl)
    // 1230 vient de la base 2018
	idParcelle := 1230
	var err error
	query ="insert into parcelle(id,id_proprietaire,code,surface,id_commune) values($1,$2,$3,$4,$5) returning id"
	_, err = db.Exec(
		query,
		idParcelle,
		1,
		"0M0168",
		58.864,
		7)
	if err != nil {
		panic(err)
	}
	query = "insert into parcelle_lieudit(id_parcelle,id_lieudit) values($1,$2)"
	_, err = db.Exec(
		query,
		idParcelle,
		358)
	if err != nil {
		panic(err)
	}
	//
	query = "insert into parcelle_fermier(id_parcelle,id_fermier) values($1,$2)"
	_, err = db.Exec(
		query,
		idParcelle,
		44)
	if err != nil {
		panic(err)
	}
	//
	// Ajout liens parcelle - ug
	//
	query = "insert into parcelle_ug(id_parcelle,id_ug) values($1,$2)"
	for _, idUG := range([]int{282,624,513,33}){
        _, err = db.Exec(
            query,
            idParcelle,
            idUG)
        if err != nil {
            panic(err)
        }
	}
	
	fmt.Println("Migration effectuée : 2024-04-17-ugs-orphelines")
}
