import React from "react";
import type { ProgressVariant } from "../../../types";
import { cn } from "../../lib/utils";

interface ProgressProps extends React.HTMLAttributes<HTMLDivElement> {
  value?: number;
  max?: number;
  variant?: ProgressVariant;
  className?: string;
  showLabel?: boolean;
}

const Progress = React.forwardRef<HTMLDivElement, ProgressProps>(
  (
    {
      className,
      value = 0,
      max = 100,
      variant = "default",
      showLabel = false,
      ...props
    },
    ref
  ) => {
    const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

    const variants: Record<ProgressVariant, string> = {
      default: "bg-gradient-to-r from-neobrutal-blue to-blue-500",
      success: "bg-gradient-to-r from-neobrutal-green to-emerald-400",
      warning: "bg-gradient-to-r from-amber-500 to-yellow-400",
      danger: "bg-gradient-to-r from-red-500 to-red-400",
    };

    return (
      <div className={cn("relative", className)}>
        <div
          ref={ref}
          className={cn(
            "w-full h-2.5 bg-gray-100 border-2 border-black rounded-sm overflow-hidden"
          )}
          {...props}
        >
          <div
            className={cn(
              "h-full transition-all duration-500 ease-out",
              variants[variant]
            )}
            style={{ width: `${percentage}%` }}
          />
        </div>
        {showLabel && (
          <span className="absolute right-0 -top-6 text-xs font-semibold text-gray-600">
            {Math.round(percentage)}%
          </span>
        )}
      </div>
    );
  }
);

Progress.displayName = "Progress";

export { Progress };
