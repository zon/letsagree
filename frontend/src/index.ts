import { parseArgs } from "node:util"
import { Elysia } from "elysia"
import { aboutPage, homePage, loginPage, logout, notHumanPage } from "./orchestration.js"
import { fetchBackendVersion, fetchUser, logout as backendLogout } from "./backend.js"

const { values } = parseArgs({
  options: {
    addr: { type: "string", default: ":3000" },
    backend: { type: "string", default: "http://localhost:8080" },
  },
})

const frontendVersion = process.env.VERSION ?? "dev"

const app = new Elysia()
  .get("/about", async (ctx) => {
    const html = await aboutPage(
      () => fetchBackendVersion(values.backend as string),
      frontendVersion,
    )
    return new Response(html, { headers: { "Content-Type": "text/html; charset=utf-8" } })
  })
  .get("/", async (ctx) => {
    const session = ctx.cookie["session"]?.value ?? null
    const response = await homePage(session, () => fetchUser(values.backend as string, session ?? ""))
    if (response.type === "redirect") {
      return Response.redirect(response.to, 302)
    }
    return new Response(response.content, { headers: { "Content-Type": "text/html; charset=utf-8" } })
  })
  .get("/login", (ctx) => {
    const session = ctx.cookie["session"]?.value ?? null
    const response = loginPage(session)
    if (response.type === "redirect") {
      return Response.redirect(response.to, 302)
    }
    return new Response(response.content, { headers: { "Content-Type": "text/html; charset=utf-8" } })
  })
  .get("/not-human", (ctx) => {
    const response = notHumanPage()
    return new Response(response.content, { headers: { "Content-Type": "text/html; charset=utf-8" } })
  })
  .post("/logout", async (ctx) => {
    const session = ctx.cookie["session"]?.value ?? null
    const response = await logout(() => backendLogout(values.backend as string, session ?? ""))
    return Response.redirect(response.to, 302)
  })

const addrValue = values.addr as string
const portMatch = addrValue.match(/^:?(\d+)$/)
const port = portMatch ? parseInt(portMatch[1], 10) : 3000

app.listen(port, () => {
  console.log(`Frontend listening on ${addrValue}`)
})