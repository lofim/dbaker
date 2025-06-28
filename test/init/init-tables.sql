create table public.users (
    id          integer primary key generated always as identity,
    first_name  varchar(255) not null,
    last_name   text not null unique,
    description text null,
    deleted_at  timestamp null,
    role        varchar(10)
);

create table public.groups (
    id integer primary key generated always as identity,
    group_name varchar(255) not null
);

create table public.users_groups (
    user_id integer,
    group_id integer,

    primary key (user_id, group_id),
    foreign key (user_id) references public.users (id),
    foreign key (group_id) references public.groups (id)
);
