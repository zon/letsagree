import { describe, test, expect } from "bun:test";
import { renderAbout } from "./render.js";

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