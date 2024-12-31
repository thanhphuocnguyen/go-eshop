CREATE TABLE
    "sessions" (
        "session_id" uuid PRIMARY KEY,
        "user_id" bigint NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
        "refresh_token" varchar NOT NULL,
        "user_agent" varchar NOT NULL,
        "client_ip" varchar NOT NULL,
        "is_blocked" boolean NOT NULL DEFAULT false,
        "expired_at" timestamptz NOT NULL,
        "created_at" timestamptz NOT NULL DEFAULT now ()
    );

CREATE INDEX ON "sessions" ("user_id");