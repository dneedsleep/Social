ALTER TABLE
    users
ADD
    COLUMN password VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE
    users
ADD
    COLUMN created_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
