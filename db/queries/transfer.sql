-- name: CreateTransfer :one
insert into transfers (from_account_id, to_account_id, amount)
values ($1, $2, $3)
returning *;

-- name: GetTransfer :one
select *
from transfers
where id = $1 limit 1;

-- name: ListTransfers :many
select *
from transfers
order by id
limit $1
    offset $2;

-- name: UpdateTransfer :one
update transfers
set amount = $2
where id = $1
returning *;

-- name: DeleteTransfer :exec
delete
from transfers
where id = $1;