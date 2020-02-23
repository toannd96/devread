-- +migrate Up
CREATE TABLE "users" (
  "user_id" text PRIMARY KEY,
  "full_name" text,
  "email" text UNIQUE,
  "password" text,
  "role" text,
  "created_at" TIMESTAMPTZ NOT NULL,
  "updated_at" TIMESTAMPTZ NOT NULL
);

CREATE TABLE "posts" (
  "name" text PRIMARY KEY,
  "url" text,
  "category" text,
  "view" text,
  "share" text,
  "comment" text,
  "build_by" text,
  "time_read" text,
  "created_at" TIMESTAMPTZ NOT NULL
);

CREATE TABLE "bookmarks" (
  "bid" text PRIMARY KEY,
  "user_id" text,
  "post_name" text,
  "created_at" TIMESTAMPTZ NOT NULL,
  "updated_at" TIMESTAMPTZ NOT NULL,
  unique (user_id, post_name)
);

ALTER TABLE "bookmarks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
ALTER TABLE "bookmarks" ADD FOREIGN KEY ("post_name") REFERENCES "posts" ("name");

-- +migrate Down
DROP TABLE bookmarks;
DROP TABLE users;
DROP TABLE posts;
