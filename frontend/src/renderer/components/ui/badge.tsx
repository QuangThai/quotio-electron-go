import React from "react";
import type { ButtonVariant } from "../../../types";
import { cn } from "../../lib/utils";

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: ButtonVariant;
  children: React.ReactNode;
}

const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({ className, variant = "default", children, ...props }, ref) => {
    const variants: Record<ButtonVariant, string> = {
      default: "bg-white text-black border-2 border-black",
      primary: "bg-neobrutal-blue text-white border-2 border-black",
      secondary: "bg-surface text-black border-2 border-black",
      success: "bg-neobrutal-green text-white border-2 border-black",
      warning: "bg-amber-400 text-black border-2 border-black",
      danger: "bg-red-500 text-white border-2 border-black",
      ghost: "bg-transparent text-black border-2 border-transparent",
      purple: "bg-neobrutal-purple text-white border-2 border-black",
    };

    return (
      <span
        ref={ref}
        className={cn(
          "inline-flex items-center gap-1 px-2 py-0.5 text-[10px] font-bold rounded-base uppercase tracking-tighter",
          "transition-colors duration-150",
          variants[variant],
          className
        )}
        {...props}
      >
        {children}
      </span>
    );
  }
);

Badge.displayName = "Badge";

export { Badge };
