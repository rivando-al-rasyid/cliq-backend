CREATE TABLE "links"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "user_id" uuid NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "origin_link" varchar UNIQUE NOT NULL,
    "slug" varchar(32) UNIQUE NOT NULL,
    "is_deleted" boolean NOT NULL DEFAULT FALSE,
    "created_at" timestamp NOT NULL DEFAULT (now()),
    "deleted_at" timestamp,
    "updated_at" timestamp DEFAULT (now())
);

CREATE INDEX idx_links_slug ON links(slug);
