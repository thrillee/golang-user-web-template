-- Create "otps" table
CREATE TABLE "public"."otps" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "target_ref" text NULL,
  "target" text NULL,
  "code" text NULL,
  "expiration" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_otps_deleted_at" to table: "otps"
CREATE INDEX "idx_otps_deleted_at" ON "public"."otps" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NULL,
  "phone" text NULL,
  "password" text NULL,
  "account_ref" text NULL,
  "full_name" bytea NULL,
  "display_picture" bytea NULL,
  "username" text NULL,
  "age" bigint NULL,
  "active" boolean NULL,
  "is_verified" boolean NULL,
  "last_login" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_account_ref" to table: "users"
CREATE UNIQUE INDEX "idx_users_account_ref" ON "public"."users" ("account_ref");
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
