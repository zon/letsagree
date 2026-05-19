import { describe, test, expect } from "bun:test";
import { renderAbout, renderHome, renderLogin, renderNotHuman } from "./render.js";
import type { User } from "./backend.js";

describe("renderAbout", () => {
  test("returns an HTML string", () => {
    const html = renderAbout("1.0.0", "2.0.0");
    expect(typeof html).toBe("string");
    expect(html.includes("<!DOCTYPE html>")).toBe(true);
  });

  test("page displays the frontendVersion string", () => {
    const html = renderAbout("1.0.0", "2.0.0");
    expect(html).toContain("1.0.0");
    expect(html).toContain("Frontend Version");
  });

  test("page displays the backendVersion string when it is not null", () => {
    const html = renderAbout("1.0.0", "2.0.0");
    expect(html).toContain("2.0.0");
    expect(html).toContain("Backend Version");
  });

  test("page displays an unavailable indicator when backendVersion is null", () => {
    const html = renderAbout("1.0.0", null);
    expect(html).toContain("unavailable");
  });

  test("page is styled with Tailwind CSS", () => {
    const html = renderAbout("1.0.0", "2.0.0");
    expect(html).toContain("tailwind");
    expect(html).toContain("class=");
  });
});

describe("renderHome", () => {
  test("returns an HTML string", () => {
    const user = { id: "user-123", email: "alice@example.com", name: "Alice" } as User;
    const html = renderHome(user);
    expect(typeof html).toBe("string");
    expect(html.includes("<!DOCTYPE html>")).toBe(true);
  });

  test("displays the user's name", () => {
    const user = { id: "user-123", email: "alice@example.com", name: "Alice" } as User;
    const html = renderHome(user);
    expect(html).toContain("Alice");
  });

  test("includes a logout form that POSTs to /logout", () => {
    const user = { id: "user-123", email: "alice@example.com", name: "Alice" } as User;
    const html = renderHome(user);
    expect(html).toContain('method="post"');
    expect(html).toContain('action="/logout"');
  });
});

describe("renderLogin", () => {
  test("returns an HTML string", () => {
    const html = renderLogin();
    expect(typeof html).toBe("string");
    expect(html.includes("<!DOCTYPE html>")).toBe(true);
  });

  test("includes a Login with Humanity Protocol button that links to GET /auth/login", () => {
    const html = renderLogin();
    expect(html).toContain("Login with Humanity Protocol");
    expect(html).toContain('href="/auth/login"');
  });
});

describe("renderNotHuman", () => {
  test("returns an HTML string", () => {
    const html = renderNotHuman();
    expect(typeof html).toBe("string");
    expect(html.includes("<!DOCTYPE html>")).toBe(true);
  });

  test("indicates that biometric verification via Humanity Protocol is required", () => {
    const html = renderNotHuman();
    expect(html).toContain("Biometric verification required");
    expect(html).toContain("Humanity Protocol");
  });
});