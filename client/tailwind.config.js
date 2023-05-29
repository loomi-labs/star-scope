const { addDynamicIconSelectors } = require('@iconify/tailwind');

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [ "./src/**/*.rs", "./index.html" ],
  theme: {
    extend: {
      colors: {
        transparent: 'transparent',
        current: 'currentColor',
        primary: '#D68940',
        "primary-gradient": {
          from: '#FA803D',
            to: '#D14809',
        },
        d: {
          bg: '#111223',
          "bg-1": '#2D1B2E',
        },
        'badge-blue': '#40FBFC',
        'badge-red': '#EF3C65',
        'badge-orange': '#F07D6C',
        'badge-green': '#11AA89',
        black: '#000000',
        white: '#FFFFFF',
        purple: {
          50: '#fff2ff',
          100: '#fdeafd',
          200: '#f2dff2',
          300: '#e0cee1',
          400: '#bba9bc',
          500: '#9a899b',
          600: '#977E98',
          700: '#342335',
          800: '#2D1B2E',
          900: '#1e111f'
        },
      }
    }
  },
  plugins: [
    // Iconify plugin
    addDynamicIconSelectors(),
  ],
  darkMode: 'class',
}

