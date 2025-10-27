CREATE TABLE problems (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    example_input TEXT,
    example_output TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE test_cases (
    id BIGSERIAL PRIMARY KEY,
    problem_id BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_sample BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language VARCHAR(50) NOT NULL,
    source_code TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'Queued',
    result_message TEXT,
    execution_time FLOAT,
    memory_used FLOAT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
