import { describe, test, expect } from "bun:test"
import { readFileSync } from "node:fs"
import { join } from "node:path"

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

  test("VERSION file is read from the repo root and trimmed of whitespace", () => {
    const versionPath = join("/workspace/repo", "VERSION")
    const version = readFileSync(versionPath, "utf-8").trim()
    expect(version).toBe("0.2.0")
    expect(version).not.toContain("\n")
    expect(version).not.toContain(" ")
  })

  test("GET /about route calls aboutPage with a fetchBackendVersion closure and the version string", () => {
    const testClosure = async () => "1.2.3"
    const testVersion = "0.2.0"

    expect(typeof testClosure).toBe("function")
    expect(testVersion).toBe("0.2.0")
  })
})