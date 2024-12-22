CREATE TABLE complaints (
    id varchar(255) unique not null primary key,
    user_id int unsigned not null,
    title varchar(255) not null,
    complaints text not null,
    response json not null,
    created_at timestamp not null,
    constraint fk_user_id_complaints foreign key (user_id) references evia_dev.users(id)
) engine innodb;