CREATE TABLE IF NOT EXISTS users (
    "user_id" TEXT not null,
    "created_at" datetime not null default CURRENT_TIMESTAMP,
    "email" TEXT not null,
    "first_name" TEXT not null,
    "last_name" TEXT not null,
    "nickname" TEXT null,
    "birthday" datetime not null,
    "image" TEXT null,
    "about" TEXT null,
    "status" TEXT not null default PUBLIC,
    "password" TEXT not null,
    primary key ("user_id")
);
