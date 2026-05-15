# Let's Agree

A website where users vote on positions to discover the most popular consensus. The [Humanity Protocol](https://www.humanity.org) ensures each real person gets exactly one vote.

## How it works

Users browse or create position statements and cast a single vote per topic. Votes are aggregated to surface the strongest points of agreement across the community. Identity is verified through Humanity Protocol, preventing duplicate or bot votes.

## Tech stack

**Backend** — Go, [Gin](https://gin-gonic.com), [GORM](https://gorm.io), [Kong](https://github.com/alecthomas/kong)

**Frontend** — [Bun](https://bun.sh), [Elysia](https://elysiajs.com), [htmx](https://htmx.org), [Tailwind CSS](https://tailwindcss.com)

**Task runner** — [just](https://github.com/casey/just)

## Development

Start the backend (with live reload):

```sh
just backend
```

Start the frontend (with live reload):

```sh
just frontend
```

Run tests:

```sh
just test
```

## License

See [LICENSE](LICENSE).
