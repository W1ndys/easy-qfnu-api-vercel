import defaultTheme from 'tailwindcss/defaultTheme'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        surface: '#E0E5EC',
        ink: '#3D4852',
        muted: '#6B7280',
        accent: '#6C63FF',
        'accent-light': '#8B84FF',
        success: '#38B2AC',
        placeholder: '#A0AEC0'
      },
      fontFamily: {
        sans: ['"DM Sans"', ...defaultTheme.fontFamily.sans],
        display: ['"Plus Jakarta Sans"', '"DM Sans"', ...defaultTheme.fontFamily.sans]
      },
      boxShadow: {
        neu: '9px 9px 16px rgba(163,177,198,0.6), -9px -9px 16px rgba(255,255,255,0.55)',
        'neu-hover': '12px 12px 20px rgba(163,177,198,0.7), -12px -12px 20px rgba(255,255,255,0.65)',
        'neu-small': '5px 5px 10px rgba(163,177,198,0.6), -5px -5px 10px rgba(255,255,255,0.55)',
        'neu-inset': 'inset 6px 6px 10px rgba(163,177,198,0.6), inset -6px -6px 10px rgba(255,255,255,0.55)',
        'neu-inset-deep': 'inset 10px 10px 20px rgba(163,177,198,0.7), inset -10px -10px 20px rgba(255,255,255,0.65)',
        'neu-inset-small': 'inset 3px 3px 6px rgba(163,177,198,0.6), inset -3px -3px 6px rgba(255,255,255,0.55)',
        'neu-accent': '8px 8px 18px rgba(92,86,219,0.28), -8px -8px 18px rgba(255,255,255,0.48)'
      },
      borderRadius: {
        neu: '32px',
        'neu-sm': '16px'
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' }
        }
      },
      animation: {
        float: 'float 3s ease-in-out infinite'
      }
    },
  },
  plugins: [],
}
