CREATE TABLE IF NOT EXISTS posts (
    "post_id" TEXT not null,
    "group_id" TEXT null,
    "created_by" TEXT not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "content" TEXT null,
    "image" TEXT null,
    "visibility" TEXT null default PUBLIC,
    primary key ("post_id")
);
