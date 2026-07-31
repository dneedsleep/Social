ALTER TABLE 
    posts
CONSTRAINT fk_user FOREIGN KEY (user_id ) REFERENCES user (id);