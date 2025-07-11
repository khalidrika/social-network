CREATE TABLE IF NOT EXISTS comments (
    "comment_id" TEXT not null,
    "post_id" TEXT not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "created_by" TEXT not null,
    "content" text null,
    "image" TEXT null,
    primary key ("comment_id")
);
