import { describe, test, expect } from "bun:test"
import {
  aboutPage,
  assertRedirectsTo,
  assertRendersHome,
  assertRendersLogin,
  assertRendersNotHuman,
  homePage,
  logout,
  PageResponse,
} from "./orchestration"
import { sessions, users, backend } from "./backend"

const backendReturning =
  (version: string) => async () => version

const backendUnavailable = () => async () => null

function assertContainsFrontendVersion(html: string, version: string) {
  expect(html).toContain(version)
  expect(html).toContain("Frontend Version")
}

function assertContainsBackendVersion(html: string, version: string) {
  expect(html).toContain(version)
  expect(html).toContain("Backend Version")
}

function assertBackendVersionUnavailable(html: string) {
  expect(html).toContain("unavailable")
}

describe("assertRedirectsTo", () => {
  test("passes when response is a redirect to the expected path", () => {
    const response: PageResponse = { type: "redirect", to: "/login" }
    expect(() => assertRedirectsTo(response, "/login")).not.toThrow()
  })

  test("throws when response type is not redirect", () => {
    const response: PageResponse = { type: "html", content: "<p>test</p>" }
    expect(() => assertRedirectsTo(response, "/login")).toThrow(
      "Expected redirect to /login but got html",
    )
  })

  test("throws when redirect target does not match", () => {
    const response: PageResponse = { type: "redirect", to: "/" }
    expect(() => assertRedirectsTo(response, "/login")).toThrow(
      "Expected redirect to /login but got redirect to /",
    )
  })
})

describe("assertRendersHome", () => {
  test("passes when response is html containing home page markers", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Welcome, Alice</p><form action=\"/logout\" method=\"POST\">",
    }
    expect(() =>
      assertRendersHome(response, { id: "1", name: "Alice", email: "alice@example.com" })
    ).not.toThrow()
  })

  test("throws when response type is not html", () => {
    const response: PageResponse = { type: "redirect", to: "/login" }
    expect(() =>
      assertRendersHome(response, { id: "1", name: "Alice", email: "alice@example.com" })
    ).toThrow("Expected html but got redirect")
  })

  test("throws when content does not include user name", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Welcome</p><form action=\"/logout\" method=\"POST\">",
    }
    expect(() =>
      assertRendersHome(response, { id: "1", name: "Alice", email: "alice@example.com" })
    ).toThrow("Home page content missing expected user name")
  })

  test("throws when content does not include logout form", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Welcome, Alice</p>",
    }
    expect(() =>
      assertRendersHome(response, { id: "1", name: "Alice", email: "alice@example.com" })
    ).toThrow("Home page content missing logout form")
  })
})

describe("assertRendersLogin", () => {
  test("passes when response is html containing login page marker", () => {
    const response: PageResponse = {
      type: "html",
      content: "<button>Login with Humanity Protocol</button>",
    }
    expect(() => assertRendersLogin(response)).not.toThrow()
  })

  test("throws when response type is not html", () => {
    const response: PageResponse = { type: "redirect", to: "/login" }
    expect(() => assertRendersLogin(response)).toThrow(
      "Expected html but got redirect",
    )
  })

  test("throws when content does not include login button", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Welcome</p>",
    }
    expect(() => assertRendersLogin(response)).toThrow(
      "Login page content missing expected marker",
    )
  })
})

describe("assertRendersNotHuman", () => {
  test("passes when response is html containing biometric error marker", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Biometric verification required</p>",
    }
    expect(() => assertRendersNotHuman(response)).not.toThrow()
  })

  test("throws when response type is not html", () => {
    const response: PageResponse = { type: "redirect", to: "/not-human" }
    expect(() => assertRendersNotHuman(response)).toThrow(
      "Expected html but got redirect",
    )
  })

  test("throws when content does not include biometric error marker", () => {
    const response: PageResponse = {
      type: "html",
      content: "<p>Error</p>",
    }
    expect(() => assertRendersNotHuman(response)).toThrow(
      "Not human page content missing expected biometric marker",
    )
  })
})

describe("homePage", () => {
  test("home page redirects to login when session cookie is absent", async () => {
    const response = await homePage(sessions.absent(), backend.userNotFound())
    assertRedirectsTo(response, "/login")
  })

  test("home page redirects to login when session is invalid", async () => {
    const response = await homePage(sessions.any(), backend.userNotFound())
    assertRedirectsTo(response, "/login")
  })

  test("home page renders for authenticated user", async () => {
    const user = users.any()
    const response = await homePage(sessions.any(), backend.findingUser(user))
    assertRendersHome(response, user)
  })
})

describe("GET /about shows both versions when backend is available", () => {
  test("about page shows both versions when backend is available", async () => {
    const html = await aboutPage(backendReturning("2.0.0"), "1.0.0")
    assertContainsFrontendVersion(html, "1.0.0")
    assertContainsBackendVersion(html, "2.0.0")
  })
})

describe("GET /about when backend is unreachable", () => {
  test("about page marks backend unavailable when backend cannot be reached", async () => {
    const html = await aboutPage(backendUnavailable(), "1.0.0")
    assertContainsFrontendVersion(html, "1.0.0")
    assertBackendVersionUnavailable(html)
  })
})

describe("logout", () => {
  test("logout redirects to login", async () => {
    const response = await logout(backend.thatLogsOut())
    assertRedirectsTo(response, "/login")
  })
})