create table if not exists submissions (
    id bigserial primary key ,
    user_id bigint not null references users(id) on delete cascade ,
    problem_id bigint not null references problems(id) on delete cascade,
    language varchar(50) not null ,
    source_code text,
    status varchar(50) not null,
    result_message text not null,
    execution_time double precision not null,
    memory_used double precision not null,
    submitted_at timestamp default CURRENT_TIMESTAMP
)