# Entain test task

All database tables are in `migrations` folder. To run them, use `make` command.

1. `make migrate-up` - to run up migrations
2. `make migrate-down` - to run down migrations
3. `make migrate-force` - to force run migrations if you have some errors like `error: Dirty database version -1. Fix and force version.`

### Tables:

1. `users` - contains users data
2. `transaction` - contains transactions data

To run the app, use `make run` command.

### Endpoints to test:

1. `GET /users` - to get all users
2. `GET /users/{user_id}` - to get user by id, check his balance
3. `GET /transactions/{user_id}` - to get all transactions by user id (check if user has any transactions)
4. `POST /process-record/{user_id}` - to process record by user id

Process record request body example:

```
{
    "amount": 10,
    "transaction_id": "64338a05-81e5-426b-b01e-927e447c9e33",
    "state": "win"
}
```

Transaction id is unique, so you can't process the same transaction twice, provide UUID v4 format.
State can be `win` or `lose`.
Amount is a number, can be positive or negative.

### Required header for all endpoints:

1. `Source-Type: game` - available values: `game`, `server`, `payment`

Postman collection is in `postman` folder to test endpoints.

### To run application in docker container:

1. `make docker-build` - to build docker image
2. `make docker-run` - to run docker container
