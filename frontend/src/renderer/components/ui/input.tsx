import React from "react";
import { cn } from "../../lib/utils";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  className?: string;
  type?: string;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type = "text", ...props }, ref) => {
    return (
      <input
        type={type}
        ref={ref}
        className={cn(
          "w-full px-3 py-2 bg-white border-2 border-black rounded-base shadow-neobrutal-sm",
          "placeholder:text-gray-400 placeholder:font-normal",
          "focus:outline-none focus:ring-2 focus:ring-neobrutal-blue/30 focus:border-neobrutal-blue",
          "focus:shadow-neobrutal focus:translate-x-[-1px] focus:translate-y-[-1px]",
          "transition-all duration-200 font-bold text-sm",
          "disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-gray-50",
          className
        )}
        {...props}
      />
    );
  }
);

Input.displayName = "Input";

export { Input };
