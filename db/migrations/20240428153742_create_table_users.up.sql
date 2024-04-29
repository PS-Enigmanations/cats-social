-- Create table users
CREATE TABLE "public"."users" (
    "id" serial NOT NULL,
    "uuid" UUID NOT NULL,
    "email" varchar(150) NOT NULL,
    "password" varchar(150) NOT NULL,
    "phone" TEXT,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    "deleted_at" timestamptz DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE UNIQUE INDEX idx_credentials_email ON "public"."users" USING btree (email);
CREATE UNIQUE INDEX idx_credentials_phone ON "public"."users" USING btree (phone);
