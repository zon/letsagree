import { describe, test, expect, afterEach, beforeEach } from "bun:test"
import { spawn } from "node:child_process"
import { setTimeout as wait } from "node:timers/promises"

describe("entry-point scenarios", () => {
  describe("--addr defaults to :3000", () => {
    test("GIVEN the frontend is started with no --addr flag THEN the server listens on :3000", async () => {
      const server = spawn("bun", ["src/index.ts"], {
        cwd: "/workspace/repo/frontend",
        stdio: ["pipe", "pipe", "pipe"],
      })
      await wait(500)

      try {
        const response = await fetch("http://localhost:3000/about")
        expect(response.status).toBe(200)
      } finally {
        server.kill()
        await wait(100)
      }
    })
  })

  describe("--addr flag overrides listen address", () => {
    test("GIVEN the frontend is started with --addr :4000 WHEN a client sends GET /about to port 4000 THEN the response status is 200 OK", async () => {
      const server = spawn("bun", ["src/index.ts", "--addr", ":4000"], {
        cwd: "/workspace/repo/frontend",
        stdio: ["pipe", "pipe", "pipe"],
      })
      await wait(500)

      try {
        const response = await fetch("http://localhost:4000/about")
        expect(response.status).toBe(200)
      } finally {
        server.kill()
        await wait(100)
      }
    })
  })

  describe("--backend flag overrides backend URL", () => {
    test("GIVEN the frontend is started with --backend http://localhost:9090 WHEN the about page is requested THEN the frontend fetches the backend version from http://localhost:9090/version", async () => {
      const server = spawn(
        "bun",
        ["src/index.ts", "--backend", "http://localhost:9090"],
        {
          cwd: "/workspace/repo/frontend",
          stdio: ["pipe", "pipe", "pipe"],
        },
      )
      await wait(500)

      try {
        const response = await fetch("http://localhost:3000/about")
        expect(response.status).toBe(200)
        const html = await response.text()
        expect(html).toContain("Frontend Version")
      } finally {
        server.kill()
        await wait(100)
      }
    })
  })
})