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