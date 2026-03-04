import * as React from "react"
import { X } from "lucide-react"

import { cn } from "@/lib/utils"

// Simple toast primitive components
export interface ToastProps extends React.HTMLAttributes<HTMLDivElement> {
    variant?: "default" | "destructive" | "success" | "warning"
}

const Toast = React.forwardRef<HTMLDivElement, ToastProps>(
    ({ className, variant = "default", ...props }, ref) => {
        const variants = {
            default: "bg-background text-foreground border-border",
            destructive: "destructive group border-destructive bg-destructive text-destructive-foreground",
            success: "border-green-500/50 bg-green-50 text-green-900 dark:bg-green-900/20 dark:text-green-400",
            warning: "border-yellow-500/50 bg-yellow-50 text-yellow-900 dark:bg-yellow-900/20 dark:text-yellow-400",
        }

        return (
            <div
                ref={ref}
                className={cn(
                    "group pointer-events-auto relative flex w-full items-center justify-between space-x-4 overflow-hidden rounded-md border p-6 pr-8 shadow-lg transition-all",
                    variants[variant],
                    className
                )}
                {...props}
            />
        )
    }
)
Toast.displayName = "Toast"

const ToastClose = React.forwardRef<
    HTMLButtonElement,
    React.ButtonHTMLAttributes<HTMLButtonElement>
>(({ className, ...props }, ref) => (
    <button
        ref={ref}
        className={cn(
            "absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100",
            className
        )}
        {...props}
    >
        <X className="h-4 w-4" />
    </button>
))
ToastClose.displayName = "ToastClose"

const ToastTitle = React.forwardRef<
    HTMLHeadingElement,
    React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
    <h3
        ref={ref}
        className={cn("text-sm font-semibold", className)}
        {...props}
    />
))
ToastTitle.displayName = "ToastTitle"

const ToastDescription = React.forwardRef<
    HTMLParagraphElement,
    React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
    <div
        ref={ref}
        className={cn("text-sm opacity-90", className)}
        {...props}
    />
))
ToastDescription.displayName = "ToastDescription"

export { Toast, ToastClose, ToastTitle, ToastDescription }
