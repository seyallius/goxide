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

    window.scrollTo(0, 0);
  } catch (error) {
    contentEl.innerHTML =
      "<h1>404 - Page Not Found</h1><p>Looks like this page wandered off... (╥_╥)</p>";
    console.error("Failed to load content:", error);
  }
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
