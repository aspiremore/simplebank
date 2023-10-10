CREATE TABLE "Users"
(
    "username"            varchar PRIMARY KEY,
    "email"               varchar UNIQUE NOT NULL,
    "created_at"          timestamptz    NOT NULL DEFAULT 'now()',
    "password_changed_at" timestamptz    NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "hashed_password"     varchar        NOT NULL,
    "full_name"           varchar        NOT NULL
);

--CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");


ALTER TABLE "accounts" ADD constraint  "owner_currency_key" unique ("owner", "currency");
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "Users" ("username");
