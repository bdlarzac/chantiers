
-- opérations simples d'un chantier plaquettes
create type typeop as enum(
    'AB', -- abattage
    'DB', -- debardage
    'BR', -- broyage
    'DC'  -- dechiquetage
);
