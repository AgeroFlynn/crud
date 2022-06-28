CREATE extension IF NOT EXISTS "uuid-ossp";

SET bytea_output = 'escape';

CREATE TABLE IF NOT EXISTS users (
                       user_id       UUID DEFAULT uuid_generate_v4 (),
                       name          TEXT,
                       email         TEXT UNIQUE,
                       roles         TEXT[],
                       password_hash bytea,
                       date_created  TIMESTAMP,
                       date_updated  TIMESTAMP,

                       PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS phone_dict (
                          phone_dict_id   UUID DEFAULT uuid_generate_v4 (),
                          telegram         TEXT,
                          user_id      UUID,
                          date_created TIMESTAMP,
                          date_updated TIMESTAMP,

                          PRIMARY KEY (phone_dict_id),
                          FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);