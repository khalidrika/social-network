CREATE TABLE IF NOT EXISTS "event" (
    "event_id" TEXT not null,
    "group_id" TEXT not null,
    "created_by" TEXT not null,
    "created_at" DATETIME not null default CURRENT_TIMESTAMP,
    "title" TEXT not null,
    "content" TEXT not null,
    "date" DATETIME not null,
    primary key ("event_id")
);