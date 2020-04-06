
create table stockloyer (
    id                      serial primary key,
    id_stockage             int not null references stockage(id),
    montant                 numeric not null,
    datedeb                 date not null,
    datefin                 date not null
);
create index stockloyer_id_stockage_idx on stockloyer(id_stockage);
