CREATE TABLE IF NOT EXISTS messages (
			"message_id" TEXT not null,
			"sender_id" TEXT not null,
			"receiver_id" TEXT not null,
			"type" TEXT not null,
			"created_at" datetime not null default CURRENT_TIMESTAMP,
			"content" TEXT not null,
			"is_read" INT default 0,
			primary key ("message_id")
		);
