
-- types de valorisation pour les chantiers "autres valorisations"
create type typevalorisation as enum(
    'BO', -- bois d'oeuvre
    'CH', -- bois de chauffage
    'PI', -- piquets
    'PL', -- palette
    'PP'  -- pâte à papier
);
                                                                                                                  