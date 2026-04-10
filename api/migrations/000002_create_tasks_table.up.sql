CREATE TYPE status AS ENUM('open','done');

CREATE TABLE IF NOT EXISTS tasks(
id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
project_id uuid,
name VARCHAR NOT NULL,
status status NOT NULL DEFAULT 'open',
start_at DATE NOT NULL,
created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
last_update TIMESTAMPTZ,
FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
);
