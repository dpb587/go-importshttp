module.exports = {
  purge: {
    mode: 'all',
    preserveHtmlElements: false,
    options: {
      keyframes: true,
    },
    content: [
      '../*.html',
      './src/inlineassets/*.svg',
    ],
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        gopherblue: '#00ADD8',
        darkblue: '#007D9C',
        darkerblue: '#006C8B',
        lightblue: '#5DC9E2',
        fuschia: '#CE3262',
        aqua: '#00A29C',
        black: '#000000',
        yellow: '#FDDD00',
      },
    },
    screens: {
      'sm': '640px',
    },
  },
  variants: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}
