module.exports = {
  content: ["./templates/*.{html,js}", "./templates/static/*.{html,js}"],
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
  extend: {
    animation: {
      'sprite': 'spriteAnimation 1s steps(5) infinite',
    },
    keyframes: {
      spriteAnimation: {
        '0%': { backgroundPosition: '0 0' },
        '100%': { backgroundPosition: '-1350px 0' },
      }
    }
  }
}
