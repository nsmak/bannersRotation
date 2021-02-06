-- +goose Up
CREATE TABLE IF NOT EXISTS banner (
    id serial NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS slot (
    id serial NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS social_dem (
    id serial NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS banner_slot (
    banner_id serial NOT NULL,
    slot_id serial NOT NULL,
    PRIMARY KEY (banner_id, slot_id),
    FOREIGN KEY (banner_id)
        REFERENCES banner (id),
    FOREIGN KEY (slot_id)
        REFERENCES slot (id)
);

CREATE TABLE IF NOT EXISTS banner_showing (
    banner_id serial NOT NULL,
    slot_id serial NOT NULL,
    social_id serial NOT NULL,
    date timestamptz NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (banner_id, slot_id, social_id, date),
    FOREIGN KEY (banner_id)
        REFERENCES banner (id),
    FOREIGN KEY (slot_id)
        REFERENCES slot (id),
    FOREIGN KEY (social_id)
        REFERENCES social_dem (id)
);

CREATE TABLE IF NOT EXISTS banner_click (
    banner_id serial NOT NULL,
    slot_id serial NOT NULL,
    social_id serial NOT NULL,
    date timestamptz NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (banner_id, slot_id, social_id, date),
    FOREIGN KEY (banner_id)
        REFERENCES banner (id),
    FOREIGN KEY (slot_id)
        REFERENCES slot (id),
    FOREIGN KEY (social_id)
        REFERENCES social_dem (id)
);

INSERT INTO banner (id, description) VALUES (1, 'Car banner'), (2, 'Shop banner'), (3, 'Food banner');
INSERT INTO slot (id, description) VALUES (1, 'Header slot'), (2, 'Footer slot'), (3, 'Left menu banner');
INSERT INTO social_dem (id, description) VALUES (1, 'Молодежь'), (2,'Старики'), (3, 'Сотрудники спецслужб');

-- +goose Down
drop table banner;
drop table slot;
drop table social_dem;
drop table banner_slot;
drop table banner_showing;
drop table banner_click;