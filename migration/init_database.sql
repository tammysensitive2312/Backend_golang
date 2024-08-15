-- auto-generated definition
create table users
(
    id         bigint auto_increment
        primary key,
    email      varchar(191) not null,
    password   longtext     not null,
    username   varchar(255) not null,
    created_at datetime(3)  null,
    updated_at datetime(3)  null,
    constraint uni_users_email
        unique (email)
);

-- auto-generated definition
create table projects
(
    id                 bigint auto_increment
        primary key,
    name               longtext    not null,
    category           longtext    not null,
    project_spend      bigint      null,
    project_variance   bigint      null,
    revenue_recognised bigint      null,
    project_started_at datetime(3) not null,
    project_ended_at   datetime    null,
    created_at         datetime(3) null,
    updated_at         datetime(3) null,
    deleted_at         datetime(3) null
);

create index idx_projects_deleted_at
    on projects (deleted_at);


-- auto-generated definition
create table user_projects
(
    project_id bigint   not null,
    user_id    bigint   not null,
    ID         bigint auto_increment
        primary key,
    created_at datetime null,
    updated_at datetime null,
    deleted_at datetime null,
    constraint user_projects_projects_id_fk
        foreign key (project_id) references projects (id),
    constraint user_projects_users_id_fk
        foreign key (user_id) references users (id)
);

create index fk_project_projects_user
    on user_projects (project_id);

