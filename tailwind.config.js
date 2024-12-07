/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["ui/template/**/*.templ"],
  daisyui: {
    themes: ["acid"],
  },
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
};
