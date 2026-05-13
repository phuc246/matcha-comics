/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './components/**/*.{js,vue,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './plugins/**/*.{js,ts}',
    './app.vue',
    './error.vue',
  ],
  theme: {
    extend: {
      colors: {
        'bg-primary': '#0D0F14',
        'bg-secondary': '#141720',
        'bg-tertiary': '#1C2030',
        'accent-primary': '#9CA764',
        'accent-hover': '#B8C878',
        'accent-gold': '#F1E8C7',
        'accent-dark': '#6B7445',
        'text-primary': '#F1E8C7',
        'text-secondary': '#A8A8B3',
        'text-muted': '#5C5C6B',
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
        heading: ['Montserrat', 'sans-serif'],
      },
      boxShadow: {
        'glow-sm': '0 0 10px rgba(156, 167, 100, 0.2)',
        'glow-md': '0 0 20px rgba(156, 167, 100, 0.35)',
        'glow-lg': '0 0 40px rgba(156, 167, 100, 0.5)',
        'dark-lg': '0 8px 48px rgba(0,0,0,0.8)',
      },
      animation: {
        shimmer: 'shimmer 1.5s infinite',
        'fade-up': 'fadeUp 0.5s ease forwards',
        'scale-in': 'scaleIn 0.3s ease forwards',
        'glow-pulse': 'glowPulse 2s ease-in-out infinite',
      },
      keyframes: {
        shimmer: {
          '0%': { backgroundPosition: '-1000px 0' },
          '100%': { backgroundPosition: '1000px 0' },
        },
        fadeUp: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.9)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
        glowPulse: {
          '0%, 100%': { boxShadow: '0 0 10px rgba(156, 167, 100, 0.2)' },
          '50%': { boxShadow: '0 0 30px rgba(156, 167, 100, 0.6)' },
        },
      },
    },
  },
  plugins: [],
}
