# Entain test task

All database tables are in `migrations` folder. To run them, use `make` command.

1. `make migrate-up` - to run up migrations
2. `make migrate-down` - to run down migrations
3. `make migrate-force` - to force run migrations if you have some errors like `error: Dirty database version -1. Fix and force version.`

### Tables:

1. `users` - contains users data
2. `transaction` - contains transactions data

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

### To run the app locally:

1. Create `.env` file in root folder and add all required variables from `.env.example` file
2. Run `make migrate-up` command to run migrations and create all tables with test user (user_id: 63e83104-b9a7-4fec-929e-9d08cae3f9b9)
3. Run `make run` command to run application
4. Take a look at `postman` folder to take collection for testing all endpoints

Test user with id `63e83104-b9a7-4fec-929e-9d08cae3f9b9` will be created automatically when you run migrations.
This user has 50 amount of his balance for testing.

### To run application in docker container:

1. `make docker-build` - to build docker image
2. `make docker-run` - to run docker container
