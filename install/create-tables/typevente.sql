
-- Types de vente d'un chantier autres valorisations
create type typevente as enum(
    'NON',  -- non spécifié
    'BSP', -- Bois sur pied
    'BDR', -- Bord de route
    'LIV'  -- Livré
);
