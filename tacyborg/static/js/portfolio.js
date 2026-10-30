(() => {
  const slides = [...document.querySelectorAll("[data-slide]")];
  if (!slides.length) return;

  const main = document.querySelector("#portfolio");
  const menu = document.querySelector(".portfolio-menu");
  const select = document.querySelector("#section-select");
  const number = document.querySelector("#slide-number");
  const title = document.querySelector("#slide-title");
  const announcement = document.querySelector("#slide-announcement");
  const languageSwitch = document.querySelector("#language-switch");
  const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
  const imageJobs = new WeakMap();
  let current = -1;
  let destination = null;
  let scrollFrame = 0;
  let resizeFrame = 0;
  let lastWidth = window.innerWidth;
  let settleTimer;

  const clamp = (index) => Math.max(0, Math.min(slides.length - 1, index));
  const indexFromHash = () => clamp(slides.findIndex((slide) => `#${slide.id}` === window.location.hash));

  function warmNearby(index) {
    for (let i = Math.max(0, index - 1); i <= Math.min(slides.length - 1, index + 2); i++) {
      slides[i].querySelectorAll("img").forEach((image) => {
        image.fetchPriority = i === index ? "high" : "low";
        if (imageJobs.has(image)) return;
        image.loading = "eager";
        // Failed downloads must not hide alt text or block navigation.
        const ready = image.decode().catch(() => {}).then(() => image.classList.add("image-ready"));
        imageJobs.set(image, ready);
      });
    }
  }

  function setActive(index, updateHash = true) {
    warmNearby(index);
    if (index === current) return;
    current = index;
    const slide = slides[index];
    select.value = slide.id;
    number.textContent = String(index + 1).padStart(2, "0");
    title.textContent = slide.dataset.title;
    announcement.textContent = `${index + 1} ${document.body.dataset.of} ${slides.length}: ${slide.dataset.title}`;
    document.title = `${slide.dataset.title} — Tacy Borg`;
    languageSwitch.href = `${languageSwitch.dataset.base}#${slide.id}`;
    if (updateHash) window.history.replaceState(null, "", `#${slide.id}`);
  }

  function positionOf(index) {
    return Math.max(0, Math.min(slides[index].offsetTop, main.scrollHeight - main.clientHeight));
  }

  function visibleIndex() {
    const center = main.scrollTop + main.clientHeight / 2;
    return clamp(slides.findIndex((slide) => slide.offsetTop + slide.offsetHeight > center));
  }

  function navigate(index, { history = true, smooth = true, focus = true } = {}) {
    index = clamp(index);
    destination = index;
    const hash = `#${slides[index].id}`;
    if (history && window.location.hash !== hash) window.history.pushState(null, "", hash);
    setActive(index, false);
    main.scrollTo({ top: positionOf(index), behavior: smooth && !reducedMotion.matches ? "smooth" : "instant" });
    if (focus) slides[index].focus({ preventScroll: true });
    clearTimeout(settleTimer);
    settleTimer = setTimeout(settleScroll, 180);
  }

  function settleScroll() {
    destination = null;
    setActive(visibleIndex());
  }

  main.addEventListener("scroll", () => {
    clearTimeout(settleTimer);
    settleTimer = setTimeout(settleScroll, 160);
    if (scrollFrame) return;
    scrollFrame = requestAnimationFrame(() => {
      scrollFrame = 0;
      const index = visibleIndex();
      warmNearby(index);
      if (destination === null && !resizeFrame) setActive(index);
    });
  }, { passive: true });
  main.addEventListener("scrollend", settleScroll);
  main.addEventListener("wheel", () => { destination = null; }, { passive: true });
  main.addEventListener("pointerdown", () => { destination = null; }, { passive: true });

  select.addEventListener("change", () => {
    menu.open = false;
    navigate(slides.findIndex((slide) => slide.id === select.value));
  });
  document.addEventListener("click", (event) => {
    const link = event.target.closest('a[href^="#"]');
    if (!link || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
    const index = slides.findIndex((slide) => `#${slide.id}` === link.getAttribute("href"));
    if (index < 0) return;
    event.preventDefault();
    navigate(index);
  });
  document.addEventListener("pointerdown", (event) => {
    if (menu.open && !menu.contains(event.target)) menu.open = false;
  });
  window.addEventListener("hashchange", () => navigate(indexFromHash(), { history: false, smooth: false }));
  window.addEventListener("popstate", () => navigate(indexFromHash(), { history: false, smooth: false }));
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && menu.open) {
      menu.open = false;
      menu.querySelector("summary").focus();
      return;
    }
    if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey || event.target.closest("input, textarea, select, button, summary, iframe, [contenteditable]")) return;
    if (window.getSelection()?.toString()) return;
    let index;
    if (event.key === "ArrowLeft") index = current - 1;
    if (event.key === "ArrowRight") index = current + 1;
    if (event.key === "Home") index = 0;
    if (event.key === "End") index = slides.length - 1;
    if (index === undefined) return;
    event.preventDefault();
    navigate(index);
  });
  window.addEventListener("resize", () => {
    if (resizeFrame || lastWidth === window.innerWidth) return;
    const index = current;
    resizeFrame = requestAnimationFrame(() => {
      resizeFrame = 0;
      lastWidth = window.innerWidth;
      navigate(index, { history: false, smooth: false, focus: false });
    });
  });
  window.addEventListener("beforeprint", () => slides.forEach((_, index) => warmNearby(index)));

  document.body.classList.add("scroll-ready");
  menu.hidden = false;
  setActive(indexFromHash(), false);
  navigate(current, { history: false, smooth: false, focus: false });
  document.fonts.load('100px "Pixelion"').then((fonts) => {
    if (fonts.length) document.body.classList.add("pixelion-available");
  }).catch(() => { /* Keep the fallback font if loading fails. */ });
})();
