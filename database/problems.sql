create table if not exists problems (
    id bigserial primary key ,
    title varchar(255) unique not null,
    description text not null,
    difficulty varchar(40) not null,
    example_input text,
    example_output text
);
