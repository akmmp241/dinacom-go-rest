CREATE TABLE sessions (
    id int unsigned not null auto_increment primary key,
    user_id int unsigned not null,
    token varchar(255),
    expires_at timestamp not null,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
) engine innodb;