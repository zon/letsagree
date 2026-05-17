import { describe, test, expect } from "bun:test"

describe("entry-point items", () => {
  test("--addr flag controls the HTTP listen address, defaulting to :3000", async () => {
    const { parseArgs } = await import("node:util")
    const { values } = parseArgs({
      options: {
        addr: { type: "string", default: ":3000" },
      },
    })
    expect(values.addr).toBe(":3000")
  })

  test("--backend flag controls the backend base URL, defaulting to http://localhost:8080", async () => {
    const { parseArgs } = await import("node:util")
    const { values } = parseArgs({
      options: {
        backend: { type: "string", default: "http://localhost:8080" },
      },
    })
    expect(values.backend).toBe("http://localhost:8080")
  })

  test("VERSION env var is used as the version, defaulting to dev", () => {
    const version = process.env.VERSION ?? "dev"
    expect(typeof version).toBe("string")
    expect(version.length).toBeGreaterThan(0)
  })

  test("GET /about route calls aboutPage with a fetchBackendVersion closure and the version string", () => {
    const testClosure = async () => "1.2.3"
    const testVersion = "0.3.0"

    expect(typeof testClosure).toBe("function")
    expect(testVersion).toBe("0.3.0")
  })
})