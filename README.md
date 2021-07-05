This repository comes with a `makefile` which contains some useful scripts for running and maintaining the application.

To run the application:

```sh
$ git clone https://github.com/bensooraj/rndmicu.git
$ cd rndmicu
$ cp template.env .env
# Update all the relevant env variables
$ make dc-up
```

Once the is up, please navigate to [localhost:8080](http://localhost:8080/) to open the GraphQL playground. The Playground provides documentation of all possible queries and mutations which the server supports.

A Postgres DB comes seeded with a few creators for sake of simplicity.

A sample query for creating audio shorts from the GraphQL Playground is detailed in the [PR](https://github.com/bensooraj/rndmicu/pull/2).

#### Some design choices:

I opted against using an ORM for handling all things database. Instead I rely on simple SQL scripts and type-safe generated golang code.

All database schema and query `*.sql` files are placed under `data/schema/sql/`. Once the changes are made, please run `make models` for generating the `go` types.

All GraphQL schema are placed under `graph/schema`. Once the changes are made, please `make graph-gen` for re-generating the `go` types.
