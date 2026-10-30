module.exports = {
  content: ["./templates/*.{html,js}", "./templates/static/*.{html,js}"],
  plugins: [
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
    function ({ addUtilities }) {
      addUtilities({
        ".text-shadow-neon": {
          "text-shadow": "0 0 5px #39ff14, 0 0 10px #39ff14, 0 0 20px #39ff14",
        },
      });
    },
  ],
  theme: {
    extend: {
      textShadow: {
        neon: "0 0 5px #39ff14, 0 0 10px #39ff14, 0 0 20px #39ff14",
      },
      colors: {
        "neon-green": "#39ff14",
      },
      fontFamily: {
        "press-start": ['"Press Start 2P"', "cursive"],
      },
    },
  },
};
