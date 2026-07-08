/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  darkMode: 'class',
  // Disable preflight to preserve existing custom CSS reset/base styles
  corePlugins: {
    preflight: false
  },
  // Use 'tw-' prefix to avoid any class name conflicts with existing CSS
  prefix: 'tw-',
  theme: {
    extend: {
      colors: {
        primary: 'hsl(var(--primary))',
        secondary: 'var(--clr-secondary)',
        accent: 'hsl(var(--accent))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        card: 'hsl(var(--card))',
        'card-foreground': 'hsl(var(--card-foreground))',
        muted: 'hsl(var(--muted))',
        'muted-foreground': 'hsl(var(--muted-foreground))',
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        success: 'hsl(var(--success))',
        warning: 'hsl(var(--warning))',
        destructive: 'hsl(var(--destructive))',
        saffron: 'hsl(var(--saffron))',
        'india-green': 'hsl(var(--india-green))',
        navy: 'hsl(var(--navy))',
        'chakra-blue': 'hsl(var(--chakra-blue))',
        
        // Legacy clr- prefixes as fallbacks
        'clr-primary': 'var(--clr-primary)',
        'clr-primary-hover': 'var(--clr-primary-hover)',
        'clr-secondary': 'var(--clr-secondary)',
        'clr-secondary-hover': 'var(--clr-secondary-hover)',
        'clr-accent': 'var(--clr-accent)',
        'clr-accent-hover': 'var(--clr-accent-hover)',
        'clr-danger': 'var(--clr-danger)',
        'clr-success': 'var(--clr-success)',
        'clr-surface': 'var(--clr-surface)',
        'clr-surface-alt': 'var(--clr-surface-alt)',
        'clr-text-main': 'var(--clr-text-main)',
        'clr-text-muted': 'var(--clr-text-muted)',
        'clr-border': 'var(--clr-border)',
        'clr-bg': 'var(--clr-bg)'
      },
      fontFamily: {
        heading: 'var(--font-heading)',
        body: 'var(--font-body)',
        display: ['Sora', 'sans-serif'],
        sans: ['Inter', 'sans-serif']
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
