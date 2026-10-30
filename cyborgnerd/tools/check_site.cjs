// NODE_PATH=/path/to/playwright/node_modules node tools/check_site.cjs
const { chromium } = require("playwright");
const assert = require("node:assert/strict");
const base = process.env.CYBORGNERD_URL || "http://127.0.0.1:8081";

(async () => {
  const browser = await chromium.launch({ headless: true });
  try {
    for (const width of [1440, 390]) {
      const context = await browser.newContext({ viewport: { width, height: 900 } });
      const page = await context.newPage();
      const errors = [];
      page.on("pageerror", (error) => errors.push(error.message));
      await page.goto(base);
      await page.waitForFunction(() => window.htmx);
      await page.evaluate(() => { window.documentMarker = "same-document"; });
      const request = page.waitForRequest((request) => new URL(request.url()).pathname === "/about" && request.headers()["hx-request"] === "true");
      await page.locator('nav a[href="/about"]').click();
      await request;
      await page.waitForFunction(() => document.querySelector("h1").textContent === "About");
      assert.equal(new URL(page.url()).pathname, "/about");
      assert.equal(await page.title(), "About · CyborgNerd");
      assert.equal(await page.evaluate(() => window.documentMarker), "same-document");
      assert.equal(await page.locator('nav a[aria-current="page"]').getAttribute("href"), "/about");

      await page.locator('nav a[href="/projects"]').click();
      await page.waitForFunction(() => document.querySelector("h1").textContent === "Projects");
      await page.goBack();
      await page.waitForFunction(() => document.querySelector("h1").textContent === "About");
      await page.goForward();
      await page.waitForFunction(() => document.querySelector("h1").textContent === "Projects");
      await page.reload();
      assert.equal(await page.locator("h1").textContent(), "Projects");
      assert.equal(await page.locator("#content").count(), 1);
      assert.equal(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth), true);
      assert.equal(await page.locator('img, iframe, [src*="portfolio"]').count(), 0);
      assert.deepEqual(errors, []);
      await context.close();
    }

    const context = await browser.newContext({ javaScriptEnabled: false });
    const page = await context.newPage();
    await page.goto(base);
    await page.locator('nav a[href="/about"]').click();
    assert.equal(await page.locator("h1").textContent(), "About");
    await page.locator('nav a[href="/projects"]').click();
    assert.equal(await page.locator("h1").textContent(), "Projects");
    await context.close();
    console.log("CyborgNerd browser checks passed: desktop/mobile, HTMX, history, direct URLs, no-JS.");
  } finally {
    await browser.close();
  }
})().catch((error) => { console.error(error); process.exitCode = 1; });
