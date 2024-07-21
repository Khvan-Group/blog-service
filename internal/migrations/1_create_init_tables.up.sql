create type blog_status as enum ('DRAFT', 'IN_REVIEW', 'ACTIVATED', 'REJECTED');

create table if not exists t_blog_categories
(
    code varchar(255) primary key not null,
    name varchar(255)             not null
);

create table if not exists t_blogs
(
    id         bigserial primary key not null,
    created_at timestamp             not null default now(),
    created_by varchar(255)          not null,
    updated_at timestamp,
    updated_by varchar(255),
    title      varchar(255)          not null,
    content    text                  not null,
    status     blog_status           not null default 'DRAFT',
    category   varchar(255)          not null,
    likes      int                   not null default 0,
    favorites  int                   not null default 0
);

create table if not exists t_users_blogs
(
    user_login varchar(255) not null,
    blog_id    bigint       not null,
    likes      boolean      not null default false,
    favorites  boolean      not null default false
);

create unique index udx_users_blogs on t_users_blogs (user_login, blog_id);

alter table t_users_blogs
    add constraint fk_users_blogs
        foreign key (blog_id) references t_blogs (id)
            on delete cascade;

alter table t_blogs
    add constraint fk_blog_categories
        foreign key (category) references t_blog_categories (code)
            on delete cascade;