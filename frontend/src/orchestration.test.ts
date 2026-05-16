import { describe, test, expect } from "bun:test"
import { aboutPage } from "./orchestration"

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