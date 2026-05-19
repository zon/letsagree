export function renderAbout(
  frontendVersion: string,
  backendVersion: string | null,
): string {
  const backendVersionHtml = backendVersion
    ? `<span class="text-green-600 font-semibold">${backendVersion}</span>`
    : `<span class="text-red-500 italic">unavailable</span>`

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>About</title>
  <script src="https://unpkg.com/htmx.org@1.10.1"></script>
  <link rel="stylesheet" href="https://unpkg.com/tailwindcss@4.1.6/dist/tailwind.min.css">
</head>
<body class="bg-gray-100 min-h-screen flex items-center justify-center">
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
    <h1 class="text-2xl font-bold text-gray-800 mb-4">About</h1>
    <div class="space-y-3">
      <p class="text-gray-600">
        <span class="font-medium">Frontend Version:</span>
        <span class="text-blue-600 font-semibold">${frontendVersion}</span>
      </p>
      <p class="text-gray-600">
        <span class="font-medium">Backend Version:</span>
        ${backendVersionHtml}
      </p>
    </div>
  </div>
</body>
</html>`
}

import type { User } from "./backend.js"

export function renderHome(user: User): string {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Home</title>
  <script src="https://unpkg.com/htmx.org@1.10.1"></script>
  <link rel="stylesheet" href="https://unpkg.com/tailwindcss@4.1.6/dist/tailwind.min.css">
</head>
<body class="bg-gray-100 min-h-screen flex items-center justify-center">
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
    <h1 class="text-2xl font-bold text-gray-800 mb-4">Welcome, ${user.name}</h1>
    <form method="post" action="/logout">
      <button type="submit" class="bg-red-500 hover:bg-red-600 text-white font-medium py-2 px-4 rounded">
        Log out
      </button>
    </form>
  </div>
</body>
</html>`
}

export function renderLogin(): string {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Login</title>
  <script src="https://unpkg.com/htmx.org@1.10.1"></script>
  <link rel="stylesheet" href="https://unpkg.com/tailwindcss@4.1.6/dist/tailwind.min.css">
</head>
<body class="bg-gray-100 min-h-screen flex items-center justify-center">
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
    <h1 class="text-2xl font-bold text-gray-800 mb-4">Login</h1>
    <a href="/auth/login" class="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded inline-block">
      Login with Humanity Protocol
    </a>
  </div>
</body>
</html>`
}

export function renderNotHuman(): string {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Verification Required</title>
  <script src="https://unpkg.com/htmx.org@1.10.1"></script>
  <link rel="stylesheet" href="https://unpkg.com/tailwindcss@4.1.6/dist/tailwind.min.css">
</head>
<body class="bg-gray-100 min-h-screen flex items-center justify-center">
  <div class="bg-white rounded-lg shadow-lg p-8 max-w-md w-full">
    <h1 class="text-2xl font-bold text-red-600 mb-4">Verification Required</h1>
    <p class="text-gray-600">Biometric verification required. Humanity Protocol is needed to access this resource.</p>
  </div>
</body>
</html>`
}