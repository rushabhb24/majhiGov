/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  // Disable preflight to preserve existing custom CSS reset/base styles
  corePlugins: {
    preflight: false
  },
  // Use 'tw-' prefix to avoid any class name conflicts with existing CSS
  prefix: 'tw-',
  theme: {
    extend: {
      colors: {
        primary: 'var(--clr-primary)',
        'primary-hover': 'var(--clr-primary-hover)',
        secondary: 'var(--clr-secondary)',
        'secondary-hover': 'var(--clr-secondary-hover)',
        accent: 'var(--clr-accent)',
        'accent-hover': 'var(--clr-accent-hover)',
        danger: 'var(--clr-danger)',
        success: 'var(--clr-success)',
        surface: 'var(--clr-surface)',
        'surface-alt': 'var(--clr-surface-alt)',
        'text-main': 'var(--clr-text-main)',
        'text-muted': 'var(--clr-text-muted)',
        border: 'var(--clr-border)',
        bg: 'var(--clr-bg)'
      },
      fontFamily: {
        heading: 'var(--font-heading)',
        body: 'var(--font-body)'
      },
      borderRadius: {
        sm: 'var(--border-radius-sm)',
        md: 'var(--border-radius-md)',
        lg: 'var(--border-radius-lg)',
        full: 'var(--border-radius-full)'
      },
      transitionTimingFunction: {
        bounce: 'cubic-bezier(0.34, 1.56, 0.64, 1)'
      }
    }
  },
  plugins: []
}
