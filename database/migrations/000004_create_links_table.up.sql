CREATE TABLE "links"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "user_id" uuid NOT NULL,
    "origin_link" varchar UNIQUE NOT NULL,
    "slug" token_type NOT NULL,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "deleted_at" timestamp DEFAULT (now()),
    "update_at" timestamp DEFAULT (now())
);

ALTER TABLE "links"
    ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") DEFERRABLE INITIALLY IMMEDIATE;

