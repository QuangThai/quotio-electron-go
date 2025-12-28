/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/renderer/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          "Inter",
          "system-ui",
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "sans-serif",
        ],
      },
      colors: {
        // Neobrutalism vibrant colors - premium palette
        "neobrutal-blue": "#2563EB",
        "neobrutal-green": "#059669",
        "neobrutal-purple": "#7C3AED",
        "neobrutal-orange": "#EA580C",
        "neobrutal-pink": "#DB2777",
        "neobrutal-yellow": "#D97706",
        "neobrutal-cyan": "#0891B2",
        "neobrutal-indigo": "#4F46E5",
        // Surface colors
        surface: "#FAFAFA",
        "surface-hover": "#F5F5F5",
      },
      borderRadius: {
        base: "5px",
      },
      boxShadow: {
        neobrutal: "4px 4px 0px 0px rgba(0,0,0,1)",
        "neobrutal-sm": "2px 2px 0px 0px rgba(0,0,0,1)",
        "neobrutal-lg": "6px 6px 0px 0px rgba(0,0,0,1)",
        "neobrutal-xl": "8px 8px 0px 0px rgba(0,0,0,1)",
        // Soft shadows for hover states
        "neobrutal-soft": "4px 4px 0px 0px rgba(0,0,0,0.8)",
        "neobrutal-glow": "0 0 20px rgba(37, 99, 235, 0.15)",
      },
      animation: {
        "fade-in": "fadeIn 0.2s ease-out",
        "slide-up": "slideUp 0.3s ease-out",
        "pulse-soft": "pulseSoft 2s ease-in-out infinite",
        shimmer: "shimmer 2s linear infinite",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0" },
          "100%": { opacity: "1" },
        },
        slideUp: {
          "0%": { opacity: "0", transform: "translateY(10px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
        pulseSoft: {
          "0%, 100%": { opacity: "1" },
          "50%": { opacity: "0.6" },
        },
        shimmer: {
          "0%": { backgroundPosition: "-200% 0" },
          "100%": { backgroundPosition: "200% 0" },
        },
      },
      transitionTimingFunction: {
        "bounce-sm": "cubic-bezier(0.34, 1.56, 0.64, 1)",
      },
    },
  },
  plugins: [],
};
