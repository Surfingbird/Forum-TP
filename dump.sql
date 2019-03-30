create extension citext;

create table project_bd.users (
  id bigserial not null primary key,
  fullname varchar(100) not null,
  nickname citext unique not null,
  email citext unique not null,
  about text null
)

create table project_bd.forums (
  posts bigint not null default 0,
  slug citext unique not null,
  threads int not null default 0,
  title varchar(100) not null,
  user_f citext not null
);

create table project_bd.threads (
  author citext not null,
  created timestamp with time zone default now(),
  forum citext not null,
  id bigserial primary key,
  message text not null,
  slug citext null,
  title varchar(100) not null,
  votes bigint default 0
);

create table project_bd.posts (
  author citext not null,
  created timestamp with time zone default now(),
  forum citext null,
  id bigserial primary key,
  isEdited boolean default false,
  message text not null,
  parent bigint null,
  path bigint[] not null,
  thread bigint
);

create table project_bd.votes (
  v_user citext not null,
  thread bigint not null,
  u_vote int not null
);