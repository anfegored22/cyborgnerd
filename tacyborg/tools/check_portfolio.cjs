// Run against a local server; see PORTFOLIO.md for the optional Playwright setup.
const { chromium } = require("playwright");
const assert = require("node:assert/strict");
const base = process.env.PORTFOLIO_URL || "http://127.0.0.1:8080";

async function atSection(page, id) {
  await page.waitForFunction((id) => {
    const main = document.querySelector("#portfolio");
    const slide = document.getElementById(id);
    const target = Math.min(slide.offsetTop, main.scrollHeight - main.clientHeight);
    return Math.abs(main.scrollTop - target) < 2 && document.querySelector("#section-select").value === id;
  }, id);
}

(async () => {
  const browser = await chromium.launch({ headless: true });
  try {
    for (const viewport of [{ width: 1440, height: 900 }, { width: 390, height: 844 }]) {
      const context = await browser.newContext({ viewport, reducedMotion: "reduce" });
      // Keep checks deterministic without depending on YouTube's availability.
      await context.route("https://www.youtube-nocookie.com/**", (route) => route.fulfill({ contentType: "text/html", body: "<p>Documentary player</p>" }));
      const page = await context.newPage();
      const errors = [];
      const backdrops = [];
      page.on("pageerror", (error) => errors.push(error.message));
      page.on("request", (request) => {
        if (request.url().includes("/portfolio/pages/")) backdrops.push(request.url());
      });
      for (const path of ["/", "/en"]) {
        await page.goto(base + path, { waitUntil: "domcontentloaded" });
        await page.locator("body.scroll-ready").waitFor();
        assert.equal(await page.locator("[data-slide]").count(), 18);
        assert.equal(await page.locator(".pdf-frame, #toggle-reading, #zoom-slide").count(), 0);
        assert.equal(await page.locator('a[href^="/static/portfolio/web/"]').count(), 0);
        assert.equal(await page.locator("figure img").evaluateAll((images) => images.every((image) => !image.closest("a"))), true);
        assert.equal(await page.locator(".documentary-player").getAttribute("src"), "https://www.youtube-nocookie.com/embed/G9VE13QZeps?start=15");
        assert.equal(await page.locator(".documentary-player").getAttribute("loading"), "lazy");
        assert.equal(await page.locator(".documentary-link a").getAttribute("href"), "https://www.youtube.com/watch?v=G9VE13QZeps&t=15s");
        assert.equal(await page.locator("#portfolio").evaluate((main) => main.scrollWidth <= main.clientWidth), true);

        await page.mouse.move(180, 400);
        await page.mouse.wheel(0, 350);
        await page.waitForFunction(() => document.querySelector("#portfolio").scrollTop > 100);
        assert.equal(await page.locator("#portfolio").evaluate((main) => main.scrollLeft), 0);

        await page.locator(".portfolio-menu summary").click();
        await page.selectOption("#section-select", "chakana");
        await atSection(page, "chakana");
        assert.ok((await page.locator("#language-switch").getAttribute("href")).endsWith("#chakana"));
        await page.reload({ waitUntil: "domcontentloaded" });
        await atSection(page, "chakana");
        await page.setViewportSize({ width: viewport.width + 50, height: viewport.height });
        await atSection(page, "chakana");

        // Artwork taps must not open a tab or navigate away.
        const image = page.locator("#chakana figure img").first();
        await image.click();
        assert.equal(context.pages().length, 1);
        assert.equal(new URL(page.url()).hash, "#chakana");

        await page.emulateMedia({ media: "print" });
        assert.equal(await page.locator(".reading-content").first().isVisible(), true);
        assert.equal(await page.locator(".documentary-link").isVisible(), true);
        await page.emulateMedia({ media: "screen" });
      }
      assert.deepEqual(backdrops, [], "The responsive page must never request PDF backdrops");
      assert.deepEqual(errors, []);
      if (process.env.PORTFOLIO_SCREENSHOTS) {
        await page.screenshot({ path: `${process.env.PORTFOLIO_SCREENSHOTS}-${viewport.width}.png` });
      }
      await context.close();
    }

    const noJS = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 1440, height: 900 } });
    await noJS.route("https://www.youtube-nocookie.com/**", (route) => route.fulfill({ contentType: "text/html", body: "" }));
    const page = await noJS.newPage();
    await page.goto(base + "/");
    assert.equal(await page.locator("#declaracion .prose").isVisible(), true);
    assert.equal(await page.locator(".documentary-link a").isVisible(), true);
    assert.equal(await page.locator(".pdf-frame").count(), 0);
    await noJS.close();
    console.log("Portfolio browser checks passed (desktop/mobile, ES/EN, navigation, images, embed, print, no-JS).");
  } finally {
    await browser.close();
  }
})().catch((error) => { console.error(error); process.exitCode = 1; });
