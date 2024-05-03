INSERT INTO "public"."cats" ("user_id", "name", "race", "sex", "age_in_month", "description", "has_matched", "created_at", "updated_at", "deleted_at")
VALUES
    (1, 'Fluffy', 'Persian', 'male', 24, 'Fluffy is a friendly and playful cat.', FALSE, NOW(), NOW(), NULL),
    (1, 'Whiskers', 'Maine Coon', 'female', 18, 'Whiskers loves to cuddle and enjoys exploring outdoors.', FALSE, NOW(), NOW(), NULL),
    (2, 'Mittens', 'Siamese', 'male', 12, 'Mittens has beautiful blue eyes and a loud purr.', FALSE, NOW(), NOW(), NULL),
    (3, 'Luna', 'Ragdoll', 'female', 36, 'Luna is a calm and gentle cat who loves attention.', FALSE, NOW(), NOW(), NULL);
