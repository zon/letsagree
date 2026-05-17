# Coding

Before writing Go or TypeScript code, read [Ralph's code guide](https://raw.githubusercontent.com/zon/ralph/refs/heads/main/docs/code.md)

# Tools

See our [tools doc](docs/tools.md) for the libraries and tools used in this project. Always use these tools — do not introduce alternatives.

# Testing

See our [testing doc](docs/testing.md) for guidance on when to write tests and how to verify commands.

# Versioning

When bumping the version, update **both** files together:
- `VERSION`
- `helm/Chart.yaml` (`appVersion` and `version`)

Always do a **patch bump** on the chart `version` field alongside any `appVersion` change.
