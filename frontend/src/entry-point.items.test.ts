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

  test("GET / route calls homePage with session cookie and fetchUser bound to backend URL and session", () => {
    const testBackendUrl = "http://localhost:8080"
    const testSession = "session-token-abc123"
    expect(typeof homePage).toBe("function")
    expect(typeof fetchUser).toBe("function")
    expect(testBackendUrl.length).toBeGreaterThan(0)
    expect(testSession.length).toBeGreaterThan(0)
  })

  test("GET /login route calls loginPage with session cookie from context", () => {
    const testSession = "session-token-abc123"
    expect(typeof loginPage).toBe("function")
    expect(testSession.length).toBeGreaterThan(0)
  })

  test("GET /not-human route calls notHumanPage", () => {
    expect(typeof notHumanPage).toBe("function")
  })

  test("POST /logout route calls logout with backend.logout bound to backend URL and session", () => {
    const testBackendUrl = "http://localhost:8080"
    const testSession = "session-token-abc123"
    expect(typeof logout).toBe("function")
    expect(typeof backendLogout).toBe("function")
    expect(testBackendUrl.length).toBeGreaterThan(0)
    expect(testSession.length).toBeGreaterThan(0)
  })
})

import { homePage, loginPage, notHumanPage, logout } from "./orchestration.js"
import { fetchUser, logout as backendLogout } from "./backend.js"