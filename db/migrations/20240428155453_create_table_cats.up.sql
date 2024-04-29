-- Create table cats
CREATE TABLE "public"."cats" (
    "id" serial NOT NULL,
    "user_id" int NOT NULL,
    "name" varchar(30),
    "race" varchar(20),
    "sex" varchar(6),
    "age_in_month" int,
    "description" varchar(200),
    "is_already_matched" boolean,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    "deleted_at" timestamptz DEFAULT now(),
    CONSTRAINT "fk_cats_user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id"),
    PRIMARY KEY ("id")
);