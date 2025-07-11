CREATE TABLE IF NOT EXISTS group_messages (
			"message_id" TEXT not null,
			"receiver_id" TEXT not null,
			"is_read" INT default 0
);

