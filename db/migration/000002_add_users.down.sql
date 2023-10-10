ALTER TABLE "accounts" IF EXISTS DROP constraint  "owner_currency_key" ;
ALTER TABLE "accounts" IF EXISTS DROP constraint  "accounts_owner_fkey" ;
DROP TABLE if EXISTS "users";