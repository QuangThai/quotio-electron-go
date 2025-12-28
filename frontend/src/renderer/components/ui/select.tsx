import React from 'react';
import { cn } from '../../lib/utils';

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  className?: string;
  children: React.ReactNode;
}

const Select = React.forwardRef<HTMLSelectElement, SelectProps>(({
  className,
  children,
  ...props
}, ref) => {
  return (
    <select
      ref={ref}
      className={cn(
        'w-full px-4 py-3 bg-white border-4 border-black rounded-none shadow-neobrutal-sm',
        'focus:outline-none focus:shadow-neobrutal focus:translate-x-[-2px] focus:translate-y-[-2px]',
        'transition-all font-medium cursor-pointer',
        className
      )}
      {...props}
    >
      {children}
    </select>
  );
});

Select.displayName = 'Select';

export { Select };

