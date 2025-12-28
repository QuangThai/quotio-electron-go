import React, { useState } from "react";
import { cn } from "../../lib/utils";

interface TabsProps {
  defaultValue: string;
  children: React.ReactNode;
  className?: string;
}

const Tabs = ({ defaultValue, children, className }: TabsProps) => {
  const [activeTab, setActiveTab] = useState(defaultValue);

  return (
    <div className={cn("", className)}>
      {React.Children.map(children, (child) => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child, { activeTab, setActiveTab });
        }
        return child;
      })}
    </div>
  );
};

interface TabsListProps {
  children: React.ReactNode;
  activeTab: string;
  setActiveTab: (tab: string) => void;
  className?: string;
}

const TabsList = ({
  children,
  activeTab,
  setActiveTab,
  className,
}: TabsListProps) => {
  return (
    <div
      className={cn(
        "flex gap-1 mb-4 flex-wrap p-1 bg-surface rounded-base border-2 border-black",
        className
      )}
    >
      {React.Children.map(children, (child) => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child, { activeTab, setActiveTab });
        }
        return child;
      })}
    </div>
  );
};

interface TabsTriggerProps {
  value: string;
  children: React.ReactNode;
  activeTab: string;
  setActiveTab: (tab: string) => void;
  className?: string;
}

const TabsTrigger = ({
  value,
  children,
  activeTab,
  setActiveTab,
  className,
}: TabsTriggerProps) => {
  const isActive = activeTab === value;

  return (
    <button
      onClick={() => setActiveTab(value)}
      className={cn(
        "px-3 py-1.5 font-bold rounded-base transition-all duration-200",
        isActive
          ? "bg-neobrutal-blue text-white border-2 border-black shadow-neobrutal-sm"
          : "bg-transparent text-black border-2 border-transparent hover:bg-white hover:border-black/10 hover:shadow-neobrutal-sm",
        className
      )}
    >
      {children}
    </button>
  );
};

interface TabsContentProps {
  value: string;
  children: React.ReactNode;
  activeTab: string;
  className?: string;
}

const TabsContent = ({
  value,
  children,
  activeTab,
  className,
}: TabsContentProps) => {
  if (activeTab !== value) return null;

  return <div className={cn("animate-fade-in", className)}>{children}</div>;
};

export { Tabs, TabsContent, TabsList, TabsTrigger };
