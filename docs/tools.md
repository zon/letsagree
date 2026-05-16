# Tools

## General

**Task Runner:** [just](https://github.com/casey/just)
- Used to run common development processes via a `justfile`

## Containers

**Container Engine:** [Podman](https://podman.io)
- Used to build and push container images

## Backend

**Language:** Go

**CLI Parsing:** [Kong](https://github.com/alecthomas/kong)
- Used in both the REST backend and CLI apps for command line argument parsing

**ORM:** [GORM](https://gorm.io)
- Database access and schema management
- Used for all model definitions, queries, and migrations

**Testing:** [Testify](https://github.com/stretchr/testify)
- Assertions and test suite utilities

**HTTP Server:** [Gin](https://gin-gonic.com)
- HTTP routing and middleware
- Run with [Air](https://github.com/air-verse/air) in development for live reloading

**API:** REST with JSON
- JSON responses via Gin context handlers

## CLI

Command line utility apps are written in Go.

## Frontend

**Runtime:** [Bun](https://bun.sh)
- JavaScript runtime and package manager for the frontend server

**CLI Parsing:** [parseArgs](https://nodejs.org/api/util.html#utilparseargsconfig)
- Built-in Node.js utility for command line argument parsing

**HTTP Server:** [Elysia](https://elysiajs.com)
- HTTP routing and middleware for the frontend server
- Uses [elysia-livereload](https://github.com/ayaoxincheng/elysia-livereload) in development for live reloading
- Run with `bun watch` in development

**Testing:** [Bun Test](https://bun.sh/docs/cli/test)
- Built-in test runner used for all frontend tests

**Schema:** [TypeBox](https://github.com/sinclairzx81/typebox)
- Runtime type validation and JSON schema definitions

**HTTP Client:** [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)
- Built-in Web API used for all HTTP requests from the frontend

**Interactivity:** [htmx](https://htmx.org)
- Declarative HTML-driven AJAX, enabling dynamic page updates without heavy JavaScript

**Styling:** [Tailwind CSS](https://tailwindcss.com)
- Utility-first CSS framework for building the UI
