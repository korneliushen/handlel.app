/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
      colors: {
        'borderColor': '#ABABAB',
        "mainPurple": "#7A38D0"
      }
    }
  },

  plugins: [require("@tailwindcss/typography")]
};
