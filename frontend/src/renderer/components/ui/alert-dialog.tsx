import { AlertTriangle, CheckCircle, Info } from "lucide-react";
import React from "react";
import type { AlertDialogProps, ButtonVariant } from "../../../types";
import { cn } from "../../lib/utils";
import { Button } from "./button";

const AlertDialog = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  variant = "danger",
  confirmText = "Confirm",
  cancelText = "Cancel",
  className,
  ...props
}: AlertDialogProps & React.HTMLAttributes<HTMLDivElement>) => {
  if (!isOpen) return null;

  const variants: Record<
    string,
    {
      icon: typeof AlertTriangle;
      iconColor: string;
      titleColor: string;
      confirmVariant: ButtonVariant;
    }
  > = {
    danger: {
      icon: AlertTriangle,
      iconColor: "text-red-600",
      titleColor: "text-red-600",
      confirmVariant: "danger",
    },
    warning: {
      icon: AlertTriangle,
      iconColor: "text-yellow-600",
      titleColor: "text-yellow-600",
      confirmVariant: "warning",
    },
    info: {
      icon: Info,
      iconColor: "text-blue-600",
      titleColor: "text-blue-600",
      confirmVariant: "primary",
    },
    success: {
      icon: CheckCircle,
      iconColor: "text-green-600",
      titleColor: "text-green-600",
      confirmVariant: "success",
    },
  };

  const currentVariant = variants[variant] || variants.danger;
  const Icon = currentVariant.icon;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={onClose}
      {...props}
    >
      <div
        className={cn(
          "bg-white border-4 border-black shadow-neobrutal-lg p-6 max-w-md w-full mx-4",
          className
        )}
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start gap-4 mb-4">
          <div className={cn("flex-shrink-0", currentVariant.iconColor)}>
            <Icon className="w-8 h-8" />
          </div>
          <div className="flex-1">
            <h3
              className={cn(
                "text-xl font-bold mb-2",
                currentVariant.titleColor
              )}
            >
              {title}
            </h3>
            <p className="text-gray-700">{message}</p>
          </div>
        </div>

        <div className="flex gap-3 justify-end">
          <Button variant="secondary" onClick={onClose}>
            {cancelText}
          </Button>
          <Button
            variant={currentVariant.confirmVariant}
            onClick={() => {
              onConfirm();
              onClose();
            }}
          >
            {confirmText}
          </Button>
        </div>
      </div>
    </div>
  );
};

export { AlertDialog };
