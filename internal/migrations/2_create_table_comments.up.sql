create table if not exists t_comments
(
    id         bigserial primary key not null,
    created_at timestamp             not null default now(),
    created_by varchar(255)          not null,
    blog_id    bigint                not null,
    comment    text                  not null
);

create unique index udx_blog_comment on t_comments (blog_id, created_by);

alter table t_comments
    add constraint fk_comments_blogs
        foreign key (blog_id) references t_blogs (id)
            on delete cascade;