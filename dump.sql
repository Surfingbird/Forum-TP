create extension citext;

SET LOCAL synchronous_commit to OFF;

create table users (
  id bigserial not null,
  fullname varchar(250) not null,
  nickname citext unique not null,
  email citext unique not null,
  about text null,
-- --
  PRIMARY KEY (id)
);
--
-- --
create table forums (
  posts   bigint        not null default 0,
  slug    citext unique not null,
  threads int           not null default 0,
  title   varchar(250)  not null,
  user_f  citext        not null,
-- --
--   FOREIGN KEY (user_f) references users (nickname) ,
  PRIMARY KEY (slug)
);
create index idx_user_on_forum on forums using btree (user_f);
-- --
create table threads (
  author citext not null,
  created timestamp with time zone default now(),
  forum citext not null,
  id bigserial not null,
  message text not null,
  slug citext null,
  title varchar(250) not null,
  votes bigint default 0,
-- --
--   FOREIGN KEY (author) references users (nickname),
--   FOREIGN KEY (forum) references forums (slug),
  PRIMARY KEY (id)
);
create index idx_thread_slug on threads using btree (slug);
create index idx_thread_all on threads (id, slug, title, forum, author, created, message, votes);

-- --
create table posts (
  author citext not null,
  created timestamp with time zone default now(),
  forum citext null,
  id bigserial not null,
  isEdited boolean default false,
  message text not null,
  parent bigint null,
  path bigint[] not null,
  thread bigint,
-- --
--   FOREIGN KEY (author) references users (nickname),
--   FOREIGN KEY (forum) references forums (slug),
--   FOREIGN KEY (thread) references threads (id),
  PRIMARY KEY (id)
);
create index idx_check_parent on posts using btree (id, thread);
-- --
create unlogged table votes (
  v_user citext not null,
  thread bigint not null,
  u_vote int not null
);
create index idx_vote_user on votes using btree (v_user, thread);
--