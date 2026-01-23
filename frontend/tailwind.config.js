/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Dark theme
        dark: {
          base: '#1a1a1a',
          surface: '#262626',
          elevated: '#333333',
          hover: '#3d3d3d',
          border: '#404040',
        },
        // Light theme
        light: {
          base: '#ffffff',
          surface: '#fafafa',
          elevated: '#f5f5f5',
          hover: '#e5e5e5',
          border: '#e5e5e5',
        },
        // Accent color
        accent: {
          DEFAULT: '#d97706',
          hover: '#b45309',
          subtle: '#78350f',
          'subtle-light': '#fef3c7',
        },
        // HTTP methods
        method: {
          get: '#22c55e',
          post: '#3b82f6',
          put: '#f59e0b',
          patch: '#8b5cf6',
          delete: '#ef4444',
          options: '#6b7280',
          head: '#6b7280',
        },
        // Status codes
        status: {
          success: '#22c55e',
          redirect: '#3b82f6',
          'client-error': '#f59e0b',
          'server-error': '#ef4444',
        },
      },
      borderRadius: {
        sm: '4px',
        md: '8px',
        lg: '12px',
      },
    },
  },
  plugins: [],
}
