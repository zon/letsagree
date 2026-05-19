import { renderAbout, renderHome, renderLogin } from "./render.js"
import type { User } from "./backend.js"

export type PageResponse =
	| { type: "html"; content: string }
	| { type: "redirect"; to: string }

export function assertRedirectsTo(response: PageResponse, path: string): void {
	if (response.type !== "redirect") {
		throw new Error(`Expected redirect to ${path} but got html`)
	}
	if (response.to !== path) {
		throw new Error(`Expected redirect to ${path} but got redirect to ${response.to}`)
	}
}

export function assertRendersHome(response: PageResponse, user: User): void {
	if (response.type !== "html") {
		throw new Error("Expected html but got redirect")
	}
	if (!response.content.includes(user.name)) {
		throw new Error("Home page content missing expected user name")
	}
	if (!response.content.includes('/logout"')) {
		throw new Error("Home page content missing logout form")
	}
}

export function assertRendersLogin(response: PageResponse): void {
	if (response.type !== "html") {
		throw new Error("Expected html but got redirect")
	}
	if (!response.content.includes("Login with Humanity Protocol")) {
		throw new Error("Login page content missing expected marker")
	}
}

export function assertRendersNotHuman(response: PageResponse): void {
	if (response.type !== "html") {
		throw new Error("Expected html but got redirect")
	}
	if (!response.content.includes("Biometric verification required")) {
		throw new Error("Not human page content missing expected biometric marker")
	}
}

export async function aboutPage(
	fetchBackendVersion: () => Promise<string | null>,
	frontendVersion: string,
): Promise<string> {
	const backendVersion = await fetchBackendVersion()
	return renderAbout(frontendVersion, backendVersion)
}

export async function homePage(
	session: string | null,
	fetchUser: () => Promise<User | null>,
): Promise<PageResponse> {
	if (!session) return { type: "redirect", to: "/login" }
	const user = await fetchUser()
	if (!user) return { type: "redirect", to: "/login" }
	return { type: "html", content: renderHome(user) }
}

export async function logout(doLogout: () => Promise<void>): Promise<PageResponse> {
	await doLogout()
	return { type: "redirect", to: "/login" }
}

export function loginPage(session: string | null): PageResponse {
	if (session) return { type: "redirect", to: "/" }
	return { type: "html", content: renderLogin() }
}