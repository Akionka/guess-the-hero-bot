DROP TABLE IF EXISTS user_answers;
DROP TABLE IF EXISTS question_options;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS match_player_items;
DROP TABLE IF EXISTS match_players;
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS player_accounts;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS heroes;
DROP TYPE IF EXISTS position_enum;
DROP TYPE IF EXISTS team_enum;
CREATE TYPE team_enum AS ENUM ('Radiant', 'Dire');

CREATE TYPE position_enum AS ENUM ('Carry', 'Mid', 'Offlane', 'Soft Support', 'Hard Support');

CREATE TABLE heroes (
    hero_id INTEGER PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    short_name VARCHAR(64) NOT NULL
);

CREATE TABLE items (
    item_id INTEGER PRIMARY KEY,
    display_name VARCHAR(255) NOT NULL,
    short_name VARCHAR(64) NOT NULL
);

CREATE TABLE player_accounts (
    player_steam_id BIGINT NOT NULL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    is_pro BOOLEAN DEFAULT FALSE NOT NULL,
    pro_name VARCHAR(256) DEFAULT '' NOT NULL
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(32) NOT NULL UNIQUE,
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    player_steam_id BIGINT REFERENCES player_accounts
);

CREATE TABLE matches (
    match_id BIGINT NOT NULL PRIMARY KEY,
    winning_team team_enum NOT NULL,
    actual_rank INTEGER NOT NULL,
    started_at TIMESTAMP NOT NULL,
    avg_mmr INT
);

CREATE TABLE match_players (
    player_steam_id BIGINT NOT NULL REFERENCES player_accounts ON UPDATE CASCADE ON DELETE CASCADE,
    match_id BIGINT NOT NULL REFERENCES matches ON UPDATE CASCADE ON DELETE CASCADE,
    hero_id INTEGER NOT NULL REFERENCES heroes ON UPDATE CASCADE ON DELETE CASCADE,
    team team_enum NOT NULL,
    position POSITION_ENUM NOT NULL,
    PRIMARY KEY (player_steam_id, match_id),
    UNIQUE (match_id, position, team)
);

CREATE TABLE match_player_items (
    player_steam_id BIGINT NOT NULL,
    match_id BIGINT NOT NULL,
    item_id INTEGER NOT NULL REFERENCES ITEMS ON UPDATE CASCADE ON DELETE CASCADE,
    "order" INTEGER NOT NULL,
    PRIMARY KEY (player_steam_id, match_id, "order"),
    FOREIGN KEY (player_steam_id, match_id) REFERENCES match_players ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE questions (
    question_id UUID PRIMARY KEY,
    match_id BIGINT NOT NULL,
    steam_player_id BIGINT NOT NULL,
    telegram_file_id VARCHAR(128) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (steam_player_id, match_id),
    FOREIGN KEY (match_id) REFERENCES matches (match_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (player_steam_id, match_id) REFERENCES match_players (player_steam_id, match_id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE question_options (
    question_id UUID NOT NULL,
    hero_id INT NOT NULL,
    is_correct BOOL NOT NULL DEFAULT FALSE,
    "order" INT NOT NULL,
    telegram_file_id VARCHAR(128),
    PRIMARY KEY (question_id, hero_id),
    FOREIGN KEY (question_id) REFERENCES questions (question_id) ON DELETE CASCADE,
    FOREIGN KEY (hero_id) REFERENCES heroes (hero_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_correct_option ON question_options (question_id)
WHERE
    is_correct = TRUE;

CREATE TABLE user_answers (
    user_answer_id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL,
    question_id UUID NOT NULL,
    hero_id INT NOT NULL,
    answered_ad TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (question_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions (question_id) ON DELETE CASCADE,
    FOREIGN KEY (question_id, hero_id) REFERENCES question_options (question_id, hero_id) ON DELETE CASCADE
);
