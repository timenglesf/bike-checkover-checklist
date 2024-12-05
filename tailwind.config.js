/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["ui/template/**/*.templ"],
  daisyui: {
    themes: ["Nord"],
  },
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
};
