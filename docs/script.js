// File: docs/script.js
/**
 * script.js
 * This file acts as the routing and rendering engine for the zero-build SPA documentation site.
 * It fetches Markdown files, parses them into HTML, and manages the UI state (theme, sidebar).
 */

// ----------------------- Configuration -----------------------
const CONFIG = {
  defaultPage: "content/home.md",
  contentSelector: "#content",
  sidebarSelector: "#sidebar",
  themeKey: "goxide-docs-theme",
};

// ----------------------- Public Functions -----------------------

/**
 * Initializes the documentation SPA.
 * Sets up event listeners, configures Markdown parser, and loads the initial page.
 */
function initDocs() {
  configureMarked();
  setupTheme();
  setupNavigation();

  // Load initial page based on URL hash, or default to home
  const hash = window.location.hash.substring(1);
  const initialLink =
    document.querySelector(`.nav-link[href="#${hash}"]`) ||
    document.querySelector(".nav-link.active");

  if (initialLink) {
    loadContent(initialLink.getAttribute("data-src"));
  } else {
    loadContent(CONFIG.defaultPage);
  }
}

/**
 * Loads a markdown file, parses it, and injects it into the DOM.
 * @param {string} path - The relative path to the markdown file.
 */
async function loadContent(path) {
  const contentEl = document.querySelector(CONFIG.contentSelector);
  contentEl.innerHTML =
    '<div class="loading">Loading documentation... (๑•̀ㅂ•́)و✧</div>';

  try {
    const response = await fetch(path);
    if (!response.ok) throw new Error("Page not found");
    const markdown = await response.text();

    // Parse and inject
    contentEl.innerHTML = marked.parse(markdown);

    // Apply syntax highlighting to code blocks
    document.querySelectorAll("pre code").forEach((block) => {
      hljs.highlightElement(block);
    });

    // Scroll to top on page change
    window.scrollTo(0, 0);
  } catch (error) {
    contentEl.innerHTML =
      "<h1>404 - Page Not Found</h1><p>Looks like this page wandered off... (╥_╥)</p>";
    console.error("Failed to load content:", error);
  }
}

// ----------------------- Private Helpers -----------------------

/**
 * Configures the Marked.js parser for GitHub Flavored Markdown.
 */
function configureMarked() {
  marked.setOptions({
    gfm: true,
    breaks: false,
    headerIds: true,
  });
}

/**
 * Sets up event listeners for sidebar navigation.
 */
function setupNavigation() {
  const sidebar = document.querySelector(CONFIG.sidebarSelector);

  sidebar.addEventListener("click", (e) => {
    const link = e.target.closest(".nav-link");
    if (!link || link.classList.contains("external")) return;

    e.preventDefault();
    const src = link.getAttribute("data-src");
    const hash = link.getAttribute("href");

    // Update active state
    document
      .querySelectorAll(".nav-link")
      .forEach((l) => l.classList.remove("active"));
    link.classList.add("active");

    // Update URL hash for shareability
    history.pushState(null, "", hash);

    // Load content
    loadContent(src);

    // Close mobile menu if open
    if (window.innerWidth <= 768) {
      sidebar.classList.remove("open");
    }
  });

  // Mobile menu toggle
  document.getElementById("menu-toggle").addEventListener("click", () => {
    sidebar.classList.toggle("open");
  });
}

/**
 * Initializes and manages the Light/Dark theme toggle.
 */
function setupTheme() {
  const themeToggle = document.getElementById("theme-toggle");
  const html = document.documentElement;
  const hljsLight = document.getElementById("hljs-light");
  const hljsDark = document.getElementById("hljs-dark");

  // Check local storage or system preference
  const savedTheme = localStorage.getItem(CONFIG.themeKey);
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const initialTheme = savedTheme || (prefersDark ? "dark" : "light");

  applyTheme(initialTheme);

  themeToggle.addEventListener("click", () => {
    const current = html.getAttribute("data-theme");
    const next = current === "light" ? "dark" : "light";
    applyTheme(next);
    localStorage.setItem(CONFIG.themeKey, next);
  });

  function applyTheme(theme) {
    html.setAttribute("data-theme", theme);
    themeToggle.textContent = theme === "light" ? "🌙" : "☀️";
    hljsLight.disabled = theme === "dark";
    hljsDark.disabled = theme === "light";
  }
}

// ----------------------- Initialization -----------------------
document.addEventListener("DOMContentLoaded", initDocs);
