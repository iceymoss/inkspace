/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'var(--color-border)',
        input: 'var(--color-border)',
        ring: 'var(--color-primary)',
        background: 'var(--color-bg-primary)',
        foreground: 'var(--color-text-primary)',
        primary: {
          DEFAULT: 'var(--color-primary)',
          hover: 'var(--color-primary-hover)',
          light: 'var(--color-primary-light)',
          foreground: 'var(--color-text-inverse)',
        },
        secondary: {
          DEFAULT: 'var(--color-secondary)',
          foreground: 'var(--color-text-primary)',
        },
        accent: {
          DEFAULT: 'var(--color-accent)',
          hover: 'var(--color-accent-hover)',
          foreground: 'var(--color-text-inverse)',
        },
        destructive: {
          DEFAULT: 'var(--color-danger)',
          foreground: 'var(--color-text-inverse)',
        },
        success: {
          DEFAULT: 'var(--color-success)',
        },
        warning: {
          DEFAULT: 'var(--color-warning)',
        },
        muted: {
          DEFAULT: 'var(--color-bg-secondary)',
          foreground: 'var(--color-text-secondary)',
        },
        card: {
          DEFAULT: 'var(--color-bg-card)',
          foreground: 'var(--color-text-primary)',
        },
        popover: {
          DEFAULT: 'var(--color-bg-card)',
          foreground: 'var(--color-text-primary)',
        },
        sidebar: {
          bg: 'var(--color-sidebar-bg)',
          text: 'var(--color-sidebar-text)',
          'text-active': 'var(--color-sidebar-text-active)',
          border: 'var(--color-sidebar-border)',
        },
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
        mono: ['var(--font-mono)'],
      },
      fontSize: {
        xs: 'var(--font-size-xs)',
        sm: 'var(--font-size-sm)',
        base: 'var(--font-size-base)',
        lg: 'var(--font-size-lg)',
        xl: 'var(--font-size-xl)',
        '2xl': 'var(--font-size-2xl)',
        '3xl': 'var(--font-size-3xl)',
      },
      lineHeight: {
        tight: 'var(--line-height-tight)',
        base: 'var(--line-height-base)',
        relaxed: 'var(--line-height-relaxed)',
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
        xl: 'var(--radius-xl)',
        full: 'var(--radius-full)',
      },
      boxShadow: {
        sm: 'var(--shadow-sm)',
        md: 'var(--shadow-md)',
        lg: 'var(--shadow-lg)',
        xl: 'var(--shadow-xl)',
      },
      spacing: {
        xs: 'var(--spacing-xs)',
        sm: 'var(--spacing-sm)',
        md: 'var(--spacing-md)',
        lg: 'var(--spacing-lg)',
        xl: 'var(--spacing-xl)',
      },
      transitionDuration: {
        fast: '150ms',
        base: '200ms',
        slow: '300ms',
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'fade-in-up': {
          from: { opacity: '0', transform: 'translateY(10px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'fade-in': 'fade-in 0.3s ease-out',
        'fade-in-up': 'fade-in-up 0.4s ease-out',
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
