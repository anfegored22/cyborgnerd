module.exports = {
  content: ["./templates/*.{html,js}", "./templates/static/*.{html,js}"],
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ]
}
