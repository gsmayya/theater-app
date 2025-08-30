/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        theater: {
          primary: '#8b5cf6',
          secondary: '#06b6d4',
          accent: '#f59e0b',
          dark: '#1f2937',
          light: '#f9fafb',
        }
      }
    },
  },
  plugins: [],
}
