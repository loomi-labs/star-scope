const { addDynamicIconSelectors } = require('@iconify/tailwind');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./src/**/*.rs", "./index.html" ],
  theme: {},
  plugins: [
    // Iconify plugin
    addDynamicIconSelectors(),
  ],
  darkMode: 'class',
}

