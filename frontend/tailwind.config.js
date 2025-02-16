/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  /*
   * 
.quicksand-<uniquifier> {
  font-family: "Quicksand", serif;
  font-optical-sizing: auto;
  font-weight: 500;
  font-style: normal;
}

   * */
  theme: {
    fontFamily: {
      sans: ['Quicksand','sans-serif'],
    },
    extend: {},
  },
  plugins: [],
}
