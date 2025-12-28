import React from "react";
import type { ButtonSize, ButtonVariant } from "../../../types";
import { cn } from "../../lib/utils";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  children: React.ReactNode;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = "default",
      size = "default",
      children,
      disabled,
      ...props
    },
    ref
  ) => {
    const variants: Record<ButtonVariant, string> = {
      default:
        "bg-white text-black border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm",
      primary:
        "bg-neobrutal-blue text-white border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-blue-700",
      secondary:
        "bg-surface text-black border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-surface-hover",
      danger:
        "bg-red-500 text-white border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-red-600",
      success:
        "bg-neobrutal-green text-white border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-emerald-700",
      warning:
        "bg-amber-400 text-black border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-amber-500",
      ghost:
        "bg-transparent text-black border-2 border-transparent hover:bg-gray-100",
      purple:
        "bg-neobrutal-purple text-white border-2 border-black shadow-neobrutal hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-neobrutal-sm hover:bg-violet-700",
    };

    const sizes: Record<ButtonSize, string> = {
      default: "px-4 py-2 text-sm font-bold",
      sm: "px-2.5 py-1 text-xs font-bold",
      lg: "px-6 py-3 text-base font-bold",
    };

    return (
      <button
        ref={ref}
        disabled={disabled}
        className={cn(
          "inline-flex items-center justify-center gap-2 rounded-base font-bold cursor-pointer",
          "transition-all duration-200",
          "active:translate-x-[4px] active:translate-y-[4px] active:shadow-none",
          "disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:translate-x-0 disabled:hover:translate-y-0 disabled:hover:shadow-neobrutal",
          variants[variant],
          sizes[size],
          className
        )}
        {...props}
      >
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";

export { Button };
