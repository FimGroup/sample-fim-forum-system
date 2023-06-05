-- --------------------------- user ---------------------------
create table public."user"
(
    user_id      varchar(200) not null
        constraint user_pk
            primary key,
    username     varchar(200),
    password     varchar(500),
    nickname     varchar(200),
    time_created bigint,
    time_updated bigint
);

alter table public."user"
    owner to admin;

create unique index user_name_index
    on public."user" (username);


-- --------------------------- posts ---------------------------
create table public.posts
(
    post_id        varchar(50)   not null
        constraint posts_pk
            primary key,
    parent_post_id varchar(50),
    post_type      integer       not null,
    author_id      varchar(50)   not null,
    forum_id       varchar(50)   not null,
    title          varchar(1000) not null,
    content        text,
    status         integer       not null,
    time_created   bigint        not null,
    time_updated   bigint
);

comment
on column public.posts.post_type is '1-topic, 2-reply';

alter table public.posts
    owner to admin;

create index posts_active_posttype_index
    on public.posts (status, post_type, forum_id, time_created);


-- --------------------------- forum ---------------------------
create table public.forums
(
    forum_id          varchar(50)   not null,
    forum_name        varchar(1000) not null,
    forum_description varchar(1000),
    parent_forum_id   varchar(50),
    status            integer       not null,
    time_created      bigint        not null,
    time_updated      bigint
);

alter table public.forums
    owner to admin;



