module.exports = {
  content: ["./src/templates/*.{html,js}", "./src/templates/static/*.{html,js}"],
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
