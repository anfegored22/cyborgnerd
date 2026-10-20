(() => {
  const slides = [...document.querySelectorAll("[data-slide]")];
  if (!slides.length) return;

  const main = document.querySelector("#portfolio");
  const controls = document.querySelector(".portfolio-controls");
  const previous = document.querySelector("#previous-slide");
  const next = document.querySelector("#next-slide");
  const select = document.querySelector("#section-select");
  const number = document.querySelector("#slide-number");
  const title = document.querySelector("#slide-title");
  const progress = document.querySelector("#slide-progress");
  const announcement = document.querySelector("#slide-announcement");
  const reading = document.querySelector("#toggle-reading");
  const zoom = document.querySelector("#zoom-slide");
  let current = -1;

  function indexFromHash() {
    return slides.findIndex((slide) => `#${slide.id}` === window.location.hash);
  }

  function show(index, focus = false) {
    index = Math.max(0, Math.min(slides.length - 1, index));
    if (index === current) return;
    const activeElement = document.activeElement;
    const focusWasInSlide = current >= 0 && slides[current].contains(activeElement);
    slides.forEach((slide, i) => {
      slide.hidden = i !== index;
      slide.setAttribute("aria-label", `${i + 1} de ${slides.length}`);
      slide.setAttribute("aria-roledescription", "diapositiva");
    });
    current = index;
    const slide = slides[index];
    slide.scrollTop = 0;
    slide.scrollLeft = 0;
    previous.disabled = index === 0;
    next.disabled = index === slides.length - 1;
    select.value = slide.id;
    number.textContent = String(index + 1).padStart(2, "0");
    title.textContent = slide.dataset.title;
    progress.style.width = `${((index + 1) / slides.length) * 100}%`;
    announcement.textContent = `${index + 1} de ${slides.length}: ${slide.dataset.title}`;
    document.title = `${slide.dataset.title} — Tacy Borg`;
    // A disabled navigation button must not retain keyboard focus.
    if (focus || focusWasInSlide || (activeElement === previous && previous.disabled) || (activeElement === next && next.disabled)) {
      slide.focus({ preventScroll: true });
    }
  }

  function navigate(index, focus = false) {
    index = Math.max(0, Math.min(slides.length - 1, index));
    const hash = `#${slides[index].id}`;
    if (window.location.hash !== hash) window.history.pushState(null, "", hash);
    show(index, focus);
  }

  previous.addEventListener("click", () => navigate(current - 1));
  next.addEventListener("click", () => navigate(current + 1));
  select.addEventListener("change", () => navigate(slides.findIndex((slide) => slide.id === select.value)));
  reading.addEventListener("click", () => {
    const enabled = document.body.classList.toggle("reading-view");
    reading.setAttribute("aria-pressed", String(enabled));
    zoom.disabled = enabled;
    slides[current].scrollTop = 0;
    slides[current].scrollLeft = 0;
    announcement.textContent = enabled ? "Vista de lectura" : "Diseño original del PDF";
  });
  zoom.addEventListener("click", () => {
    const enabled = document.body.classList.toggle("zoomed");
    zoom.setAttribute("aria-pressed", String(enabled));
    zoom.innerHTML = enabled ? 'Ajustar <span aria-hidden="true">−</span>' : 'Ampliar <span aria-hidden="true">＋</span>';
    slides[current].scrollTop = 0;
    slides[current].scrollLeft = 0;
    announcement.textContent = enabled ? "Página ampliada. Desplázate para ver sus detalles." : "Página ajustada a la pantalla";
  });
  document.addEventListener("click", (event) => {
    const link = event.target.closest('a[href^="#"]');
    if (!link || event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
    const index = slides.findIndex((slide) => `#${slide.id}` === link.getAttribute("href"));
    if (index < 0) return;
    event.preventDefault();
    navigate(index, true);
  });
  let mouseDown = null;
  main.addEventListener("pointerdown", (event) => {
    mouseDown = event.pointerType === "mouse" ? { x: event.clientX, y: event.clientY } : null;
  });
  main.addEventListener("click", (event) => {
    const start = mouseDown;
    mouseDown = null;
    // Advance only on a plain mouse click, not a drag, touch-generated click,
    // text selection, interactive element, or click while reading/panning.
    if (!start || event.pointerType !== "mouse" || event.button !== 0 || event.detail !== 1 || event.defaultPrevented) return;
    if (event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
    if (document.body.classList.contains("zoomed") || document.body.classList.contains("reading-view")) return;
    if (!event.target.closest(".pdf-frame") || event.target.closest(".pdf-run, a, button, input, select, textarea, [contenteditable]")) return;
    if (Math.hypot(event.clientX - start.x, event.clientY - start.y) > 8 || window.getSelection()?.toString()) return;
    navigate(current + 1, true);
  });
  window.addEventListener("hashchange", () => show(indexFromHash(), true));
  window.addEventListener("popstate", () => show(indexFromHash(), true));
  document.addEventListener("keydown", (event) => {
    if (event.altKey || event.ctrlKey || event.metaKey || event.shiftKey || event.target.closest("input, textarea, select, button, [contenteditable]")) return;
    if (window.getSelection()?.toString()) return;
    if (document.body.classList.contains("zoomed") && !document.body.classList.contains("reading-view") && (event.key === "ArrowLeft" || event.key === "ArrowRight")) return;
    let index;
    if (event.key === "ArrowLeft") index = current - 1;
    if (event.key === "ArrowRight") index = current + 1;
    if (event.key === "Home") index = 0;
    if (event.key === "End") index = slides.length - 1;
    if (index === undefined) return;
    event.preventDefault();
    navigate(index, true);
  });

  let touchStart = null;
  main.addEventListener("touchstart", (event) => {
    // Horizontal gestures pan the enlarged page; do not change slides.
    if (document.body.classList.contains("zoomed") && !document.body.classList.contains("reading-view")) {
      touchStart = null;
      return;
    }
    touchStart = event.touches.length === 1 ? { x: event.touches[0].clientX, y: event.touches[0].clientY } : null;
  }, { passive: true });
  main.addEventListener("touchend", (event) => {
    if (!touchStart || !event.changedTouches.length) return;
    const dx = event.changedTouches[0].clientX - touchStart.x;
    const dy = event.changedTouches[0].clientY - touchStart.y;
    touchStart = null;
    if (Math.abs(dx) > 65 && Math.abs(dx) > Math.abs(dy) * 1.8 && !window.getSelection()?.toString()) {
      navigate(current + (dx < 0 ? 1 : -1), true);
    }
  }, { passive: true });
  main.addEventListener("touchcancel", () => { touchStart = null; }, { passive: true });

  show(indexFromHash());
  document.body.classList.add("carousel-ready");
  controls.hidden = false;
  // Use the supplied webfont with the PDF's original sizing once it loads.
  // Keep the fitted substitute if the font request fails.
  document.fonts.load('100px "Pixelion"').then((fonts) => {
    if (fonts.length) document.body.classList.add("pixelion-available");
  }).catch(() => { /* The bundled, fitted Silkscreen substitute remains active. */ });
})();
