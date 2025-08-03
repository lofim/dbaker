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

-- Table for supported numeric types
create table public.numbers_test (
    id serial primary key,
    small smallint,
    normal int4,
    big bigint,
    real_col real,
    double_col double precision
);

-- Table for supported special types
create table public.special_test (
    id serial primary key,
    uuid_col uuid,
    bool_col boolean
);

-- Table for supported date/time types
create table public.datetime_test (
    id serial primary key,
    date_col date,
    time_col time,
    ts_col timestamp,
    tstz_col timestamptz
);
