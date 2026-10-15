document.addEventListener("DOMContentLoaded", function () {
  var sidebar = document.querySelector(".app-sidebar");
  var toggle = document.querySelector(".mobile-toggle");

  // Keep selectable tooltip content inside the viewport; bulma-tooltip supplies
  // the trigger and arrow styling. Escape dismisses without moving focus.
  var tooltips = document.querySelectorAll(".app-tooltip[data-tooltip]");
  function positionTooltip(trigger) {
    var rect = trigger.getBoundingClientRect();
    var bubble = trigger.querySelector(".app-tooltip-content").getBoundingClientRect();
    var width = bubble.width;
    var height = bubble.height;
    var left = rect.left + rect.width / 2 - width / 2;
    var clampedLeft = Math.max(12, Math.min(left, document.documentElement.clientWidth - width - 12));
    trigger.style.setProperty("--tooltip-shift", (clampedLeft - left) + "px");
    trigger.classList.toggle("has-tooltip-bottom", rect.bottom + height + 12 <= window.innerHeight || rect.top < height + 12);
  }

  tooltips.forEach(function (trigger) {
    function revealTooltip() {
      trigger.classList.remove("is-tooltip-dismissed");
      positionTooltip(trigger);
    }
    trigger.addEventListener("pointerenter", revealTooltip);
    trigger.addEventListener("focusin", revealTooltip);
    trigger.querySelector("button").addEventListener("click", revealTooltip);
  });

  document.addEventListener("keydown", function (event) {
    if (event.key !== "Escape") return;
    tooltips.forEach(function (trigger) {
      trigger.classList.add("is-tooltip-dismissed");
    });
  });

  function repositionTooltips() {
    tooltips.forEach(function (trigger) {
      if (trigger.matches(":hover, :focus-within")) positionTooltip(trigger);
    });
  }
  window.addEventListener("resize", repositionTooltips);
  window.addEventListener("scroll", repositionTooltips, { passive: true });

  // Mobile sidebar toggle
  if (toggle && sidebar) {
    toggle.addEventListener("click", function () {
      sidebar.classList.toggle("is-open");
    });
  }

  // Theme dropdown (top-right): light / dark are stored and set as
  // data-theme on <html>; "system" clears both so prefers-color-scheme
  // rules apply. Initial application happens in an inline <head> script
  // to avoid a flash of the wrong theme; here we only sync the UI and
  // handle clicks.
  var themeControl = document.querySelector(".theme-control");
  if (themeControl) {
    var themeToggle = themeControl.querySelector(".theme-toggle");
    var themeMenu = themeControl.querySelector(".theme-menu");
    var themeOptions = themeMenu.querySelectorAll("button[data-theme-value]");
    var themeIcons = { light: "☀", system: "◑", dark: "☽" };

    function applyTheme(mode) {
      if (mode === "light" || mode === "dark") {
        document.documentElement.setAttribute("data-theme", mode);
      } else {
        document.documentElement.removeAttribute("data-theme");
      }
      themeToggle.textContent = themeIcons[mode];
      themeOptions.forEach(function (btn) {
        btn.classList.toggle("is-active", btn.getAttribute("data-theme-value") === mode);
      });
    }

    function closeThemeMenu() {
      themeMenu.hidden = true;
      themeToggle.setAttribute("aria-expanded", "false");
    }

    var savedTheme;
    try {
      savedTheme = localStorage.getItem("apicat-theme");
    } catch (e) {}
    applyTheme(savedTheme === "light" || savedTheme === "dark" ? savedTheme : "system");

    themeToggle.addEventListener("click", function () {
      themeMenu.hidden = !themeMenu.hidden;
      themeToggle.setAttribute("aria-expanded", String(!themeMenu.hidden));
    });

    themeOptions.forEach(function (btn) {
      btn.addEventListener("click", function () {
        var mode = this.getAttribute("data-theme-value");
        try {
          if (mode === "system") localStorage.removeItem("apicat-theme");
          else localStorage.setItem("apicat-theme", mode);
        } catch (e) {}
        applyTheme(mode);
        closeThemeMenu();
      });
    });

    document.addEventListener("click", function (e) {
      if (!themeMenu.hidden && !themeControl.contains(e.target)) closeThemeMenu();
    });

    document.addEventListener("keydown", function (e) {
      if (e.key === "Escape") closeThemeMenu();
    });
  }

  // Native details provide keyboard controls and work without JavaScript.
  // Remember explicit choices across pages in this tab; search expansion is temporary.
  var sidebarGroups = document.querySelectorAll(".app-sidebar details[data-group-key]");
  var sidebarStateKey = "apicat-sidebar:" + document.title;
  var sidebarState = {};
  var searchOpenState = null;
  var searchInput = document.getElementById("sidebar-search-input");
  try {
    var savedGroups = JSON.parse(sessionStorage.getItem(sidebarStateKey));
    if (savedGroups && typeof savedGroups === "object" && !Array.isArray(savedGroups)) sidebarState = savedGroups;
  } catch (e) {}

  sidebarGroups.forEach(function (group) {
    var key = group.getAttribute("data-group-key");
    if (typeof sidebarState[key] === "boolean") group.open = sidebarState[key];
    group.querySelector("summary").addEventListener("click", function () {
      if (searchInput && searchInput.value.trim()) return;
      // The native click action toggles open after this listener runs.
      sidebarState[key] = !group.open;
      try { sessionStorage.setItem(sidebarStateKey, JSON.stringify(sidebarState)); } catch (e) {}
    });
  });

  function expandSidebarParents(link) {
    var group = link.closest("details[data-group-key]");
    while (group) {
      group.open = true;
      if (searchOpenState) searchOpenState.set(group, true);
      group = group.parentElement.closest("details[data-group-key]");
    }
  }

  if (searchInput) {
    searchInput.addEventListener("input", function () {
      var query = this.value.trim().toLowerCase();
      if (query && !searchOpenState) {
        searchOpenState = new Map();
        sidebarGroups.forEach(function (group) { searchOpenState.set(group, group.open); });
      }
      document.querySelectorAll(".sidebar-group").forEach(function (group) {
        var visible = 0;
        group.querySelectorAll(".sidebar-link[data-search]").forEach(function (link) {
          var match = !query || link.getAttribute("data-search").toLowerCase().indexOf(query) !== -1;
          link.style.display = match ? "" : "none";
          if (match) visible++;
        });
        group.hidden = visible === 0;
        if (query && visible) group.open = true;
      });
      var componentsGroup = document.querySelector(".sidebar-component-groups");
      if (componentsGroup) {
        componentsGroup.hidden = !componentsGroup.querySelector(".sidebar-group:not([hidden])");
        if (query && !componentsGroup.hidden) componentsGroup.open = true;
      }
      if (!query && searchOpenState) {
        searchOpenState.forEach(function (open, group) { group.open = open; });
        searchOpenState = null;
      }
    });
  }

  // Keep nested collapse states when an ancestor is closed and reopened.
  document.querySelectorAll(".schema-table").forEach(function (table) {
    table.addEventListener("click", function (event) {
      var button = event.target.closest(".schema-toggle");
      if (!button || !table.contains(button)) return;
      button.setAttribute("aria-expanded", String(button.getAttribute("aria-expanded") !== "true"));
      var collapsedDepth = null;
      table.querySelectorAll("tbody tr").forEach(function (row) {
        var depth = Number(row.getAttribute("data-depth"));
        if (collapsedDepth !== null && depth <= collapsedDepth) collapsedDepth = null;
        row.hidden = collapsedDepth !== null;
        var toggle = row.querySelector(".schema-toggle");
        if (!row.hidden && toggle && toggle.getAttribute("aria-expanded") === "false") {
          collapsedDepth = depth;
        }
      });
    });
  });

  // Keep the selected definition visible even in long component lists.
  var activeDefinition = document.querySelector(".sidebar-component-link.is-active");
  if (activeDefinition) {
    expandSidebarParents(activeDefinition);
    activeDefinition.scrollIntoView({ block: "nearest" });
  }

  // Endpoint detail switching
  var endpointLinks = document.querySelectorAll('.sidebar-link[data-anchor]');
  if (endpointLinks.length === 0) return;

  function showEndpoint(anchor) {
    // Hide welcome message
    var welcome = document.getElementById("endpoints-welcome");
    if (welcome) welcome.style.display = "none";

    // Hide all endpoint details
    document.querySelectorAll(".endpoint-detail").forEach(function (el) {
      el.style.display = "none";
    });

    // Show selected
    var target = document.getElementById("detail-" + anchor);
    if (target) {
      target.style.display = "block";
    }

    // Update active state in sidebar
    endpointLinks.forEach(function (link) {
      link.classList.remove("is-active");
    });
    var activeLink = document.querySelector('.sidebar-link[data-anchor="' + anchor + '"]');
    if (activeLink) {
      activeLink.classList.add("is-active");
      expandSidebarParents(activeLink);
    }

    // Close mobile sidebar
    if (sidebar) {
      sidebar.classList.remove("is-open");
    }

    // Start reading the newly selected endpoint from the top
    window.scrollTo(0, 0);
  }

  endpointLinks.forEach(function (link) {
    link.addEventListener("click", function (e) {
      // Endpoint details only exist on /endpoints; elsewhere let the
      // browser navigate there and the initial-hash handler takes over.
      if (window.location.pathname !== "/endpoints") return;
      e.preventDefault();
      var anchor = this.getAttribute("data-anchor");
      showEndpoint(anchor);
      history.pushState(null, "", "/endpoints#" + anchor);
    });
  });

  // Handle initial hash
  if (window.location.pathname === "/endpoints" && window.location.hash) {
    showEndpoint(window.location.hash.slice(1));
  }

  // Handle back/forward
  window.addEventListener("popstate", function () {
    if (window.location.pathname === "/endpoints" && window.location.hash) {
      showEndpoint(window.location.hash.slice(1));
    }
  });
});
