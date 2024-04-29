-- Create table sessions
CREATE TABLE "public"."sessions" (
    "token" TEXT,
    "expires_at" timestamptz NOT NULL,
    "user_id" int NOT NULL,
     "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NULL,
	"deleted_at" timestamptz NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY ("token"),
   	CONSTRAINT fk_credentials_sessions FOREIGN KEY (user_id) REFERENCES "public"."users" (id)
);
