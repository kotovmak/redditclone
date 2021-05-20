CREATE TABLE "users" (
  "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  "username" VARCHAR NOT NULL UNIQUE,
  "encrypted_password" VARCHAR NOT NULL
);
CREATE TABLE sessions (
  "user_id" VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  "request_id" VARCHAR NOT NULL,
  "token" VARCHAR NOT NULL,
  CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE TABLE "posts" (
  "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  "score" bigint NOT NULL,
  "views" bigint NOT NULL,
  "upvote_percentage" bigint NOT NULL,
  "type" VARCHAR NOT NULL,
  "title" VARCHAR NOT NULL,
  "author" VARCHAR NOT NULL,
  "category" VARCHAR NOT NULL,
  "text" VARCHAR NULL NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  CONSTRAINT "author" FOREIGN KEY ("author") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE TABLE "comments" (
  "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  "body" VARCHAR NOT NULL,
  "author" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "post_id" VARCHAR NOT NULL,
  CONSTRAINT "post_id" FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE,
  CONSTRAINT "author" FOREIGN KEY ("author") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE TABLE "votes" (
  "id" VARCHAR NOT NULL UNIQUE PRIMARY KEY,
  "user" VARCHAR NOT NULL,
  "vote" BIGINT NOT NULL,
  "post_id" VARCHAR NOT NULL,
  CONSTRAINT "post_id" FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE,
  CONSTRAINT "user" FOREIGN KEY ("user") REFERENCES "users" ("id") ON DELETE CASCADE
);
CREATE INDEX ON "sessions" ("token");
CREATE INDEX ON "comments" ("post_id");
CREATE INDEX ON "votes" ("post_id");