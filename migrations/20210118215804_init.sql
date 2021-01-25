-- +goose Up
CREATE TABLE IF NOT EXISTS banner (
    id numeric(5) NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS slot (
    id numeric(5) NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS social_dem (
    id numeric(5) NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS banner_slot (
    banner_id numeric(5) NOT NULL,
    slot_id numeric(5) NOT NULL,
    PRIMARY KEY (banner_id, slot_id),
    FOREIGN KEY (banner_id)
        REFERENCES banner (id),
    FOREIGN KEY (slot_id)
        REFERENCES slot (id)
);

CREATE TABLE IF NOT EXISTS statistics (
    banner_id numeric(5) NOT NULL,
    slot_id numeric(5) NOT NULL,
    social_id numeric(5) NOT NULL,
    view_count integer NOT NULL DEFAULT 0,
    click_count integer NOT NULL DEFAULT 0,
    PRIMARY KEY (banner_id, slot_id, social_id),
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
drop table statistics;