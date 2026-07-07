# GoLang Booking Engine API

A GoLang API project implementing a small, vertical slice of functionality for a commercial airline flight booking system, in order to demonstrate the application of concepts from Domain-Driven Design (DDD) in GoLang, along with declarative, BDD-style tests.

## Quick-start

Requires Go, and Docker, in addition to the `make` and `curl` utilities. To start the local dependencies, run:

```
make start-local-dependencies
```

This will launch the required Postgres DB in a Docker container.

To subsequently launch the API, run:

```
make start-api
```

Curl scripts for each endpoint can be found in the `debug` folder.

## Architecture

This project uses a Hexagonal architecture in which the core business rules are implemented by a centralized Domain library, on which all other supporting libraries (database layers, Web API projects, external services) depend. This dependency direction is achieved via Inversion of Control, so that the library containing the pure business logic remains completely unaware of the technical infrastructure which supports it.

In this case, the supporting projects include a single Web API and a single access layer for a Postgres DB:

```mermaid
flowchart BT
    A[Domain]
    B[API]
    C[Postgres]

    B --> A
    B --> C
    C --> A
```

## Module structure

### Object types

**Entities** represent the real-world objects applicable to the domain being implemented — bookings and flights in this case. They implement the logic corresponding to every action which can be taken for or by the given object within the business context.

**Commands** and **queries** represent the user access model. Commands define the actions that may be taken for or on behalf of any user (often involving one or more entities), and queries define the data that may be retrieved for said users.

In this arrangement, my approach — based on the advice in the writings of Vaughn Vernon — differs from many typical DDD implementations, in which entities are used data models while the command handlers implement the core business logic. This is a common anti-pattern known as domain model anemia, which my approach has been designed to avoid. As such, the entity acts as a write-only model; it does not interact with the repository layer until its functionality is called upon. The query handlers therefore bypass the entities entirely, returning their own read-model instead.

### Folders

The entire solution is maintained as a single Go module, but is separated into three distinct folders: **domain**, **api**, and **postgres**. Packages within the latter two folders depend on packages within the `domain` folder, while those within the `domain` folder are self-contained.

#### Domain

This dependency direction is maintained by placing the contracts — the models and interfaces representing the repository layer — in `domain/contracts`, which the supporting `postgres` packages implement. Adjacent to this, entities, commands, and queries can be found in their own respective folders.

#### Postgres

Contains a Postgres implementation of the contracts in `domain`. Any type of implementation can be used and added adjacent to the `postgres` implementation; it is the `api` package that decides which one to use. This implementation connects to the Postgres DB in Docker, initializing the database with the script in the `migration` folder.

Adjacent to this is `respositories` and `queryhandlers`, containing the implementations for the repositories and query handlers defined respectively in `domain/contracts` and `domain/queryhandlers`.

#### API

Implements a Web API using the functionality from `domain` in conjunction with the infrastructure in `postgres`. This API package is responsible for bundling the `postgres` implementation together with the `domain` to create a functioning runtime. Any alternate service layer — a gRPC server, a server-side web app, a desktop app, etc — can be implemented in the same way using this type of initialization.

## Testing

For the purpose of this project, end-to-end tests have been ommitted, and the `api` and `postgres` layers are intended as simple adapters to enable debugging and demonstration; this is not intended to represent a production-ready system in its current state.

The BDD-style tests (found in `domain/commands/pencil_test.go`) cover the public interface of the domain — the command handlers — without testing the internals (the entities). This enables the entities to be refactored and replaced flexibly while the command handler tests ensure that the domain's workflow remains stable. Since the `query` section of the domain only defines a contract which the repository layer implements, query testing has been omitted as well.
