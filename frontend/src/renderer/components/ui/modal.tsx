import { X } from "lucide-react";
import type { ModalProps } from "../../../types";
import { cn } from "../../lib/utils";
import { Button } from "./button";

const Modal = ({ isOpen, onClose, title, children, className }: ModalProps) => {
  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm animate-fade-in"
      onClick={onClose}
    >
      <div
        className={cn(
          "bg-white border-2 border-black shadow-neobrutal-lg p-6 max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto rounded-sm",
          "animate-slide-up",
          className
        )}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-6">
          {title && (
            <h2 className="text-xl font-bold tracking-tight">{title}</h2>
          )}
          <Button
            variant="ghost"
            size="sm"
            onClick={onClose}
            className="ml-auto -mr-2 -mt-2 hover:bg-gray-100"
          >
            <X className="w-5 h-5" />
          </Button>
        </div>
        {children}
      </div>
    </div>
  );
};

export { Modal };
