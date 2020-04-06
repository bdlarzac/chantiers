
create table ug (
    id                      serial primary key,
    code                    varchar(10) not null, -- len max observée = 8
    type_coupe              varchar(50) not null, -- len max observée = 15 pour un seul type
    previsionnel_coupe      varchar(20) not null, -- len max observée = 13
    type_peuplement         varchar(50) not null  -- len max observée = 26
);
