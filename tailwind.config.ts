/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],

  theme: {
    extend: {
      colors: {
        'borderColor': '#ABABAB'
      }
    }
  },

  plugins: [require("@tailwindcss/typography")]
};
