export async function fetchBackendVersion(
  backendUrl: string,
): Promise<string | null> {
  try {
    const response = await fetch(`${backendUrl}/version`)
    if (!response.ok) {
      return null
    }
    return await response.text()
  } catch {
    return null
  }
}