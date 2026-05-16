import { renderAbout } from "./render.js"

export async function aboutPage(
  fetchBackendVersion: () => Promise<string | null>,
  frontendVersion: string,
): Promise<string> {
  const backendVersion = await fetchBackendVersion()
  return renderAbout(frontendVersion, backendVersion)
}