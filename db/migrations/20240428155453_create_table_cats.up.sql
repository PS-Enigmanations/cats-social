-- Create table cats
CREATE TABLE "public"."cats" (
    "id" serial NOT NULL,
    "user_id" int NOT NULL,
    "name" varchar(30),
    "race" varchar(20),
    "sex" varchar(6),
    "age_in_month" int,
    "description" varchar(200),
    "has_matched" boolean DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NULL,
	"deleted_at" timestamptz NULL,
    CONSTRAINT "fk_cats_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id"),
    PRIMARY KEY ("id")
);
