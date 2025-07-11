CREATE TABLE IF NOT EXISTS groups (
    "group_id" TEXT not null,
    "administrator" TEXT not null,
    "name" TEXT not null,
    "description" TEXT null,
    primary key ("group_id")
);
