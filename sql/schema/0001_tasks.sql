-- +goose Up
create table tasks (
    id uuid primary key default gen_random_uuid(),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    title text not null,
    description text not null,
    completed_at timestamptz null,
    due_at timestamptz null,
    deleted_at timestamptz null
);

create index idx_tasks_completed_at on tasks (completed_at);
create index idx_tasks_due_at on tasks (due_at);
create index idx_tasks_deleted_at on tasks (deleted_at);
create index idx_tasks_id_updated_at on tasks (id, updated_at);

-- +goose Down
drop table tasks;