CREATE TABLE IF NOT EXISTS "sessions" (
    "session_id" TEXT NOT NULL PRIMARY KEY,
    "user_id" TEXT NOT NULL,
    "expiration_time" DATETIME NOT NULL
);

