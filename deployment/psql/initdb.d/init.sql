CREATE TABLE IF NOT EXISTS public.decks
(
    deck_id   UUID PRIMARY KEY,
    shuffled  BOOLEAN NOT NULL DEFAULT false,
    remaining INTEGER NOT NULL DEFAULT 52
);

CREATE TABLE IF NOT EXISTS public.cards
(
    card_id SERIAL PRIMARY KEY,
    code    VARCHAR(3)  NOT NULL,
    value   VARCHAR(10) NOT NULL,
    suit    VARCHAR(10) NOT NULL,
    drawn   BOOLEAN     NOT NULL DEFAULT false,
    deck    UUID        NOT NULL,

    CONSTRAINT fk_card_deck
        FOREIGN KEY (deck)
            REFERENCES decks (deck_id)
            ON DELETE CASCADE
            DEFERRABLE INITIALLY DEFERRED
);

DROP DATABASE IF EXISTS lucky_test;
CREATE DATABASE lucky_test WITH TEMPLATE lucky OWNER db_admin;