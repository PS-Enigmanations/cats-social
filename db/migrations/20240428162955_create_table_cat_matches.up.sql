-- Create table cat matches
CREATE TABLE "public"."cat_matches" (
    "id" serial NOT NULL,
    "issued_by" int NOT NULL,
    "match_cat_id" int NOT NULL,
    "user_cat_id" int NOT NULL,
    "message" varchar(120),
    "status" varchar(10) DEFAULT 'pending',
    CONSTRAINT "fk_cat_matches_issued_by" FOREIGN KEY ("issued_by") REFERENCES "public"."users"("id"),
    CONSTRAINT "fk_cat_matches_user_cat_id" FOREIGN KEY ("user_cat_id") REFERENCES "public"."cats"("id"),
    CONSTRAINT "fk_cat_matches_match_cat_id" FOREIGN KEY ("match_cat_id") REFERENCES "public"."cats"("id"),
    PRIMARY KEY ("id")
);