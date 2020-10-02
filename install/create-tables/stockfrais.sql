
create table stockfrais (
    id                      serial primary key,
    id_stockage             int not null references stockage(id),
    typefrais               typestockfrais not null,
    montant                 numeric not null,
    datedeb                 date not null,
    datefin                 date not null,
    notes                   text
);
create index stockfrais_id_stockage_idx on stockfrais(id_stockage);
