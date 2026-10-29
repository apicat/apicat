# apicat brand assets

Official, unmodified PNGs from `apicat-guidelines-v1/assets`, supplied with the
2026-10-29 brand and development guidelines (v1.0).

- Use the matching light/dark wordmark; preserve aspect ratio and transparent padding.
- Use the dedicated 16px/32px icons for favicons and the 64px icons for compact UI.
- Do not trace, recolor, crop, or replace the logo with fonts, CSS, or library icons.
- Brand ink: `#111827`; inverse: `#F8FAFC`; tail/accent: `#FF9818`.
- Theme and semantic color tokens live in `web/src/index.css`. Keep status colors
  separate from brand orange, and use dark text on orange buttons.

Vite imports these files into the embedded `/assets/` bundle. Keep regenerated
`web/dist/` assets with changes, so the standalone Go build remains self-contained.
