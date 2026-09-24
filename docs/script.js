/**
 * script.js
 * Routing, rendering, and theme engine for the Goxide zero-build SPA.
 * Handles hash-based routing to support both sidebar and in-content anchor links.
 */

// ----------------------- Configuration -----------------------
const CONFIG = {
  defaultPage: "content/home.md",
  contentSelector: "#content",
  sidebarSelector: "#sidebar",
  themeKey: "goxide-docs-theme",
  themes: ["light", "dark", "neon"],
  themeIcons: { light: "☀️", dark: "🌙", neon: "⚡" },
  tocSelector: "#toc-list",
  tocContainer: "#toc",
};

// ----------------------- Public Functions -----------------------

/**
 * Initializes the documentation SPA.
 */
function initDocs() {
  configureMarked();
  setupTheme();
  setupNavigation();

  // Listen for URL hash changes (handles browser back/forward AND in-content links)
  window.addEventListener("hashchange", handleRouteChange);

  // Handle initial load
  handleRouteChange();
}

/**
 * Loads a markdown file, parses it, and injects it into the DOM.
 */
async function loadContent(path) {
  const contentEl = document.querySelector(CONFIG.contentSelector);
  contentEl.innerHTML =
    '<div class="loading">Loading documentation... (๑•̀ㅂ•́)و✧</div>';

  try {
    const response = await fetch(path);
    if (!response.ok) throw new Error("Page not found");
    const markdown = await response.text();

    contentEl.innerHTML = marked.parse(markdown);

    document.querySelectorAll("pre code").forEach((block) => {
      hljs.highlightElement(block);
    });

    generateTOC();

    window.scrollTo(0, 0);
  } catch (error) {
    contentEl.innerHTML =
      "<h1>404 - Page Not Found</h1><p>Looks like this page wandered off... (╥_╥)</p>";
    console.error("Failed to load content:", error);
    document.querySelector(CONFIG.tocSelector).innerHTML = "";
  }
}

// ----------------------- TOC Generation & Scroll Spy -----------------------

/**
 * Scans the content for H2/H3 headers and builds the TOC sidebar.
 */
function generateTOC() {
  const contentEl = document.querySelector(CONFIG.contentSelector);
  const tocList = document.querySelector(CONFIG.tocSelector);
  const headers = contentEl.querySelectorAll("h2, h3");

  tocList.innerHTML = ""; // Clear existing

  if (headers.length === 0) {
    document.querySelector(CONFIG.tocContainer).style.display = "none";
    return;
  } else {
    // Reset display in case it was hidden (CSS media query handles visibility mostly,
    // but we ensure inline style doesn't block it on desktop)
    document.querySelector(CONFIG.tocContainer).style.display = "";
  }

  headers.forEach((header) => {
    // Ensure header has an ID for linking
    if (!header.id) {
      header.id = header.textContent
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/(^-|-$)/g, "");
    }

    const li = document.createElement("li");
    const a = document.createElement("a");
    a.href = `#${header.id}`;
    a.textContent = header.textContent;
    a.className = header.tagName.toLowerCase(); // 'h2' or 'h3' for CSS indentation

    // Click handler for smooth scroll and hash update
    a.addEventListener("click", (e) => {
      e.preventDefault();
      const targetId = header.id;
      history.pushState(
        null,
        "",
        `#${window.location.hash.split("#")[1]}?section=${targetId}`,
      );
      // Note: Simple hash update might conflict with page routing.
      // For SPA, we usually just scroll.
      document.getElementById(targetId).scrollIntoView({ behavior: "smooth" });

      // Manually set active state immediately
      document
        .querySelectorAll("#toc-list a")
        .forEach((l) => l.classList.remove("active"));
      a.classList.add("active");
    });

    li.appendChild(a);
    tocList.appendChild(li);
  });

  // Initialize Scroll Spy
  setupScrollSpy(headers);
}

/**
 * Highlights the TOC link corresponding to the currently visible section.
 */
function setupScrollSpy(headers) {
  // Remove existing listener if any (simple approach: clone node to remove listeners?
  // Or just rely on the fact that we regenerate TOC on every page load, so old listeners die with old DOM)

  const observerOptions = {
    root: null,
    rootMargin: "-60px 0px -80% 0px", // Trigger when header is near top
    threshold: 0,
  };

  const observer = new IntersectionObserver((entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        const id = entry.target.id;
        const tocLink = document.querySelector(`#toc-list a[href="#${id}"]`);

        if (tocLink) {
          document
            .querySelectorAll("#toc-list a")
            .forEach((l) => l.classList.remove("active"));
          tocLink.classList.add("active");

          // Optional: Scroll TOC sidebar to keep active link in view
          // tocLink.scrollIntoView({ behavior: "smooth", block: "nearest" });
        }
      }
    });
  }, observerOptions);

  headers.forEach((header) => observer.observe(header));
}

// ----------------------- Routing Engine -----------------------

/**
 * Reads the current URL hash and loads the corresponding content.
 */
function handleRouteChange() {
  const hash = window.location.hash.substring(1);
  const link = document.querySelector(`.nav-link[href="#${hash}"]`);

  if (link && !link.classList.contains("external")) {
    // Update active state
    document
      .querySelectorAll(".nav-link")
      .forEach((l) => l.classList.remove("active"));
    link.classList.add("active");

    // Load content
    loadContent(link.getAttribute("data-src"));

    // Close mobile menu if open
    if (window.innerWidth <= 768) {
      document.querySelector(CONFIG.sidebarSelector).classList.remove("open");
    }
  } else if (!hash) {
    // No hash in URL, default to home
    const defaultLink =
      document.querySelector('.nav-link[href="#home"]') ||
      document.querySelector(".nav-link");
    if (defaultLink) {
      document
        .querySelectorAll(".nav-link")
        .forEach((l) => l.classList.remove("active"));
      defaultLink.classList.add("active");
      loadContent(defaultLink.getAttribute("data-src"));
    }
  }
}

// ----------------------- Private Helpers -----------------------

function configureMarked() {
  marked.setOptions({
    gfm: true,
    breaks: false,
    headerIds: true,
  });
}

function setupNavigation() {
  const sidebar = document.querySelector(CONFIG.sidebarSelector);

  // Sidebar clicks just update the URL hash. The hashchange event does the rest!
  sidebar.addEventListener("click", (e) => {
    const link = e.target.closest(".nav-link");
    if (!link || link.classList.contains("external")) return;

    e.preventDefault();
    const hash = link.getAttribute("href");

    // Update URL hash -> triggers hashchange -> loads content
    history.pushState(null, "", hash);
    handleRouteChange();
  });

  // Mobile menu toggle
  document.getElementById("menu-toggle").addEventListener("click", () => {
    sidebar.classList.toggle("open");
  });
}

/**
 * Cycles through Light -> Dark -> Neon themes.
 */
function setupTheme() {
  const themeToggle = document.getElementById("theme-toggle");
  const html = document.documentElement;
  const hljsLight = document.getElementById("hljs-light");
  const hljsDark = document.getElementById("hljs-dark");

  const savedTheme = localStorage.getItem(CONFIG.themeKey);
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const initialTheme = savedTheme || (prefersDark ? "dark" : "light");

  applyTheme(initialTheme);

  themeToggle.addEventListener("click", () => {
    const current = html.getAttribute("data-theme");
    const currentIndex = CONFIG.themes.indexOf(current);
    const nextIndex = (currentIndex + 1) % CONFIG.themes.length;
    const nextTheme = CONFIG.themes[nextIndex];

    applyTheme(nextTheme);
    localStorage.setItem(CONFIG.themeKey, nextTheme);
  });

  function applyTheme(theme) {
    html.setAttribute("data-theme", theme);
    themeToggle.textContent = CONFIG.themeIcons[theme];

    // Light uses light hljs, Dark and Neon use dark hljs (with neon CSS overrides)
    hljsLight.disabled = theme !== "light";
    hljsDark.disabled = theme === "light";
  }
}

// ----------------------- Initialization -----------------------
document.addEventListener("DOMContentLoaded", initDocs);
