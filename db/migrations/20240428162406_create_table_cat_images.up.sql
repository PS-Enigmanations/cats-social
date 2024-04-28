-- Create table cat images
CREATE TABLE "public"."cat_images" (
    "id" serial NOT NULL,
    "cat_id" int NOT NULL,
    "url" varchar(100),
    CONSTRAINT "fk_cat_images_cat_id" FOREIGN KEY ("cat_id") REFERENCES "public"."cats"("id"),
    PRIMARY KEY ("id")
);