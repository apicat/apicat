document.addEventListener("DOMContentLoaded", function () {
  var sidebar = document.querySelector(".app-sidebar");
  var toggle = document.querySelector(".mobile-toggle");

  // Mobile sidebar toggle
  if (toggle && sidebar) {
    toggle.addEventListener("click", function () {
      sidebar.classList.toggle("is-open");
    });
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
    if (window.location.hash) {
      showEndpoint(window.location.hash.slice(1));
    }
  });
});
