/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {
      colors: {
        'borderColor': '#ABABAB',
        "mainPurple": "#9516F9"
      }
    }
  },

  plugins: [require("@tailwindcss/typography")]
};
