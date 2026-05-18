export interface User {
  id: string
  email: string
  name: string
}

export async function fetchUser(
  backendUrl: string,
  session: string,
  fetchFn: typeof fetch = fetch,
): Promise<User | null> {
  try {
    const response = await fetchFn(`${backendUrl}/auth/user`, {
      headers: { Cookie: session },
    })
    if (!response.ok) {
      return null
    }
    return await response.json()
  } catch {
    return null
  }
}

export async function logout(
  backendUrl: string,
  session: string,
  fetchFn: typeof fetch = fetch,
): Promise<void> {
  await fetchFn(`${backendUrl}/auth/logout`, {
    method: "POST",
    headers: { Cookie: session },
  })
}

export const sessions = {
  any: (): string => "session-token-abc123",
  absent: (): null => null,
}

export const users = {
  any: (): User => ({
    id: "user-123",
    email: "alice@example.com",
    name: "Alice",
  }),
}

export const backend = {
  findingUser: (user: User) => async () => user,
  userNotFound: () => async () => null,
  thatLogsOut: () => async () => {},
}