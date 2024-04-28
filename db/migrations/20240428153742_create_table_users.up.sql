-- Create table users
CREATE TABLE "public"."users" (
    "id" serial NOT NULL,
    "email" varchar(150) NOT NULL,
    "name" varchar(50),
    "password" varchar(150) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Create index on email column
CREATE INDEX idx_users_email ON "public"."users" ("email");

-- Enforce uniqueness on email column
ALTER TABLE "public"."users" ADD CONSTRAINT unique_email UNIQUE ("email");