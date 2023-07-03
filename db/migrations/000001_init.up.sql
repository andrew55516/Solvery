CREATE TABLE "users"
(
    name   varchar             NOT NULL,
    class  varchar             NOT NULL,
    email  varchar PRIMARY KEY NOT NULL,
    credit int                 NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries"
(
    id         bigserial PRIMARY KEY,
    user_email varchar     NOT NULL,
    amount     int         NOT NULL,
    comment    varchar     NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries"
    ADD FOREIGN KEY (user_email) REFERENCES "users" (email);