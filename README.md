# Notes Taking Application

This repository contains the source code for a notes taking system built using Golang. The system is responsible for user sign-up and login. Once the user login is successfull, a session-id will be returned, using the session id a user will creating a note, deleting a note and will fetch all the notes.

## Prerequisites

Before running the Message Queueing System, make sure you have the following prerequisites installed on your system:

- Go programming language (go1.20.4)
- PostgreSQL(14.8)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/ankitchahal20/notes-taking-application.git
   ```

2. Navigate to the project directory:

   ```bash
   cd notes-taking-application
   ```

3. Install the required dependencies:

   ```bash
   go mod tidy
   ```

4. DB setup
    ```
    Use the scripts inside sql-scripts directory to create the tables in your db.
    ```
5. Defaults.toml
Add the values to defaults.toml and execute `go run main.go` from the cmd directory.

## APIs
There are five API's which this repo currently has.

User Sign API
```
curl -i -k -X POST \
   http://localhost:8080/v1/signup \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"name":"Ankit Chahal",
"email": "abcd11@gmail.com",
"password": "12345"
}'
```

User Login API

```
curl -i -k -X POST \
  http://localhost:8080/v1/login \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"email": "abcd11@gmail.com",
"password": "12345"
}'
```

Notes Creation API

```
curl -i -k -X POST \
  http://localhost:8080/v1/notes \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
  "sid":"dfe65d49-85ce-46ef-91f5-ab345f595164",
  "note":"hell1o"
}'
```

Get Notes API

```
curl -i -k -X GET \
  http://localhost:8080/v1/notes \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"sid": "d66aee46-d1e6-4ecd-bdab-25a334d312c3"
}'
```

Delete Note API

```
curl -i -k -X DELETE \
  http://localhost:8080/v1/notes \
  -H "transaction-id: 288a59c1-b826-42f7-a3cd-bf2911a5c351" \
  -H "content-type: application/json" \
  -d '{
"sid": "d66aee46-d1e6-4ecd-bdab-25a334d312c2",
"id": "3"
}'
```

## Project Structure

The project follows a standard Go project structure:

- `config/`: Configuration file for the application.
- `internal/`: Contains the internal packages and modules of the application.
  - `config/`: Global configuration which can be used anywhere in the application.
  - `constants/`: Contains constant values used throughout the application.
  - `db/`: Contains the database package for interacting with PostgreSQL.
  - `middleware`: Contains the logic to validate the incoming request
  - `models/`: Contains the data models used in the application.
  - `noteserror`: Defines the errors in the application
  - `service/`: Contains the business logic and services of the application.
  - `server/`: Contains the server logic of the application.
  - `utils/`: Contains utility functions and helpers.
- `main.go`: Main entry point of the application.
- `README.md`: README.md contains the description for the notes-taking-application.

## Contributing

Contributions to the Notes Taking Application are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.

## License

The Message Queueing System is open-source and released under the [MIT License](LICENSE).

## Contact

For any inquiries or questions, please contact:

- Ankit Chahal
- ankitchahal20@gmail.com

Feel free to reach out with any feedback or concerns.
