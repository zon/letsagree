import { describe, test, expect, mock } from "bun:test"
import {
  User,
  fetchUser,
  logout,
  sessions,
  users,
  backend,
} from "./backend.js"

describe("User type", () => {
  test("users.any() returns a User with arbitrary but valid fields", () => {
    const user = users.any()
    expect(typeof user.id).toBe("string")
    expect(user.id.length).toBeGreaterThan(0)
    expect(typeof user.email).toBe("string")
    expect(typeof user.name).toBe("string")
  })
})

describe("sessions", () => {
  test("sessions.any() returns a non-empty session cookie string", () => {
    const session = sessions.any()
    expect(typeof session).toBe("string")
    expect(session.length).toBeGreaterThan(0)
  })

  test("sessions.absent() returns null", () => {
    const session = sessions.absent()
    expect(session).toBeNull()
  })
})

describe("fetchUser", () => {
  const backendUrl = "http://localhost:8080"

  test("returns User when backend responds with user data", async () => {
    const mockResponse = {
      ok: true,
      json: async () => ({ id: "user-123", email: "test@example.com", name: "Test User" }),
    } as Response
    const fetchMock = mock(() => Promise.resolve(mockResponse))

    const user = await fetchUser(backendUrl, sessions.any(), fetchMock as unknown as typeof fetch)
    expect(user).not.toBeNull()
    expect(user!.id).toBe("user-123")
    expect(user!.email).toBe("test@example.com")
    expect(user!.name).toBe("Test User")
  })

  test("returns null on 401", async () => {
    const mockResponse = {
      ok: false,
      status: 401,
    } as Response
    const fetchMock = mock(() => Promise.resolve(mockResponse))

    const user = await fetchUser(backendUrl, sessions.any(), fetchMock as unknown as typeof fetch)
    expect(user).toBeNull()
  })

  test("returns null on network error", async () => {
    const fetchMock = mock(() => Promise.reject(new Error("Network error")))

    const user = await fetchUser(backendUrl, sessions.any(), fetchMock as unknown as typeof fetch)
    expect(user).toBeNull()
  })
})

describe("logout", () => {
  const backendUrl = "http://localhost:8080"

  test("sends POST /auth/logout with session as Cookie header", async () => {
    const mockResponse = {
      ok: true,
    } as Response
    let capturedRequest: Request | null = null
    const fetchMock = mock((url: string | URL | Request, options?: RequestInit) => {
      if (typeof url === "string" && url.includes("/auth/logout")) {
        capturedRequest = new Request(url, options)
      }
      return Promise.resolve(mockResponse)
    })

    await logout(backendUrl, sessions.any(), fetchMock as unknown as typeof fetch)
    expect(capturedRequest).not.toBeNull()
    expect(capturedRequest!.method).toBe("POST")
    expect(capturedRequest!.headers.get("Cookie")).toBe(sessions.any())
  })

  test("resolves without error on success", async () => {
    const mockResponse = { ok: true } as Response
    const fetchMock = mock(() => Promise.resolve(mockResponse))

    await expect(logout(backendUrl, sessions.any(), fetchMock as unknown as typeof fetch)).resolves.toBeUndefined()
  })
})

describe("backend helpers", () => {
  test("backend.findingUser(user) returns a fetchUser stub that resolves to user", async () => {
    const user = users.any()
    const stub = backend.findingUser(user)
    const result = await stub()
    expect(result).toEqual(user)
  })

  test("backend.userNotFound() returns a fetchUser stub that resolves to null", async () => {
    const stub = backend.userNotFound()
    const result = await stub()
    expect(result).toBeNull()
  })

  test("backend.thatLogsOut() returns a doLogout stub that resolves without error", async () => {
    const stub = backend.thatLogsOut()
    await expect(stub()).resolves.toBeUndefined()
  })
})