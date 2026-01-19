CREATE TABLE task_messages (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    msg_type VARCHAR(50) NOT NULL,
    body JSONB NOT NULL,
    error_details JSONB,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);