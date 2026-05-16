# Frontend Bootstrap

## Purpose

Coordinates the about page response by fetching the backend version and delegating to rendering.

## Orchestration

```typescript
// Module: frontend/src/orchestration.ts
async function aboutPage(
  fetchBackendVersion: () => Promise<string | null>,
  frontendVersion: string,
): Promise<string> {
  const backendVersion = await fetchBackendVersion()
  return renderAbout(frontendVersion, backendVersion)
}
```

`backend.fetchBackendVersion(backendUrl)` — returns the backend version string, or `null` if the backend is unavailable.  
`render.renderAbout(frontendVersion, backendVersion)` — returns an HTML string; renders an unavailable indicator when `backendVersion` is `null`.

## Tests

```typescript
// Module: frontend/src/orchestration.test.ts

test("about page shows both versions when backend is available", async () => {
  const html = await aboutPage(backendReturning("2.0.0"), "1.0.0")
  assertContainsFrontendVersion(html, "1.0.0")
  assertContainsBackendVersion(html, "2.0.0")
})

test("about page marks backend unavailable when backend cannot be reached", async () => {
  const html = await aboutPage(backendUnavailable(), "1.0.0")
  assertContainsFrontendVersion(html, "1.0.0")
  assertBackendVersionUnavailable(html)
})
```
