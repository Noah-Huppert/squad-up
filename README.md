# Squad Up
Helps plan events with no worry

# Development
Requirements:
- [Docker](https://docker.com)
- [GNU Make](https://www.gnu.org/software/make/)

## Setup
Run `make db-create` to create and start the Postgres database Docker container.

You only need to run the above command once. After than you can use 
`make db-start` and `make db-stop` to control the database container.  

If you want to destroy the database container than run `make db-destroy`.
