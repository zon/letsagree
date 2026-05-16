import { parseArgs } from "node:util"
import { readFileSync } from "node:fs"
import { join } from "node:path"
import { Elysia } from "elysia"
import { aboutPage } from "./orchestration.js"
import { fetchBackendVersion } from "./backend.js"

const { values } = parseArgs({
  options: {
    addr: { type: "string", default: ":3000" },
    backend: { type: "string", default: "http://localhost:8080" },
  },
})

const versionPath = join(process.cwd(), "..", "VERSION")
const frontendVersion = readFileSync(versionPath, "utf-8").trim()

const app = new Elysia().get("/about", async () => {
  return aboutPage(
    () => fetchBackendVersion(values.backend as string),
    frontendVersion,
  )
})

const addrValue = values.addr as string
const portMatch = addrValue.match(/^:?(\d+)$/)
const port = portMatch ? parseInt(portMatch[1], 10) : 3000

app.listen(port, () => {
  console.log(`Frontend listening on ${addrValue}`)
})