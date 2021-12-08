module.exports = {
  purge: {
    enabled: true,
    content: [
      './templates/*.html',
      './assets/*.js',
    ],
  },
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
