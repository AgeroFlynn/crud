CREATE extension IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       user_id       UUID DEFAULT uuid_generate_v4 (),
                       name          TEXT,
                       email         TEXT UNIQUE,
                       roles         TEXT[],
                       password_hash TEXT,
                       date_created  TIMESTAMP,
                       date_updated  TIMESTAMP,

                       PRIMARY KEY (user_id)
);

CREATE TABLE phone_dict (
                          phone_dict_id   UUID DEFAULT uuid_generate_v4 (),
                          telegram         TEXT,
                          user_id      UUID,
                          date_created TIMESTAMP,
                          date_updated TIMESTAMP,

                          PRIMARY KEY (phone_dict_id),
                          FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);