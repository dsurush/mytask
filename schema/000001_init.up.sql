create table tasks
(
    id serial not null unique,
    title text not null,
    description text not null,
    done boolean not null
);
