
create table lieudit_mot (
    mot                     varchar(255) not null,
    id                      int not null,
    nom                     varchar(255) not null
);
create index lieudit_mot_mot_idx on lieudit_mot(mot);