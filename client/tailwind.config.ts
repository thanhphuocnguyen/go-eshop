import type { Config } from 'tailwindcss';

export default {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        satoshi: ['Satoshi', 'sans-serif'],
        inter: ['Inter', 'sans-serif'],
      },
      colors: {
        'secondary-green': '#1b3a1a',
        'primary-green': '#3b6a3a',
        'primary-gray': '#7d93b3',
        'secondary-gray': '#8fa7bb',
        'primary-blue': '#C6E7FF',
        'secondary-blue': '#D4F6FF',
        'primary-orange': '#FF5722',
        'secondary-orange': '#FF8303',
        background: 'var(--background)',
        foreground: 'var(--foreground)',
      },
    },
  },
  plugins: [],
} satisfies Config;
