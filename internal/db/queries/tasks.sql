-- name: CreateTask :one
insert into tasks (title, description, due_at)
values ($1, $2, $3)
returning 
    id,
    created_at,
    updated_at,
    title,
    description,
    completed_at,
    due_at,
    deleted_at;

-- name: GetTask :one
select 
    id,
    created_at,
    updated_at,
    title,
    description,
    completed_at,
    due_at,
    deleted_at
from tasks
where id = $1
limit 1;

-- name: ListTasks :many
select   
    id,
    created_at,
    updated_at,
    title,
    description,
    completed_at,
    due_at,
    deleted_at
from tasks
order by created_at desc;

-- name: UpdateTask :one
update tasks
set title = $2,
    description = $3,
    due_at = $4,
    completed_at = $5,
    updated_at = now()
where id = $1
returning
    id,
    created_at,
    updated_at,
    title,
    description,
    completed_at,
    due_at,
    deleted_at;

-- name: DeleteTask :one
update tasks set
    deleted_at = now(),
    updated_at = now()
where id = $1
returning
    id,
    created_at,
    updated_at,
    title,
    description,
    completed_at,
    due_at,
    deleted_at;