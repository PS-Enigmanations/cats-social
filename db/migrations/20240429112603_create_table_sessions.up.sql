-- Create table sessions
CREATE TABLE "public"."sessions" (
    "access_token" TEXT,
    "token_expires" INT NOT NULL,
    "user_id" int NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "deleted_at" timestamptz DEFAULT now(),
    CONSTRAINT sessions_pkey PRIMARY KEY ("access_token"),
   	CONSTRAINT fk_credentials_sessions FOREIGN KEY (user_id) REFERENCES "public"."users" (id)
);
