# Backend app writed on GO and Postgres for Asperitas Reddit clone
https://github.com/d11z/asperitas
https://asperitas.vercel.app


## Installation

### Prerequisites
- go 1.16
- docker
- golang-migrate cli

1. Clone this repository

2. Install dependencies
    ```bash
    $ go mod download
    $ docker-compose up
    $ migrate -path migrations -database "postgresql://user:password@localhost:5432/reddit?sslmode=disable" up
    $ migrate -path migrations -database "postgresql://user:password@localhost:5433/reddit_test?sslmode=disable" up
    ```

## Run the app

1. Start server locally
    ```bash
    $ make
    ```

## Testing

1. Start server tests
    ```bash
    $ make test
    ```

## License

This project is made available under the **MIT License**.