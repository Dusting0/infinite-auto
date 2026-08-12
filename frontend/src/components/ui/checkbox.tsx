import type { InputHTMLAttributes, ReactNode } from "react";
import { cn } from "../../lib/utils";

interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  label: ReactNode;
}

export function Checkbox({ className, label, ...props }: CheckboxProps) {
  return (
    <label className={cn("inline-flex cursor-pointer items-center gap-2 rounded-md border border-border bg-background/80 px-3 py-2 text-sm text-foreground transition-colors hover:bg-accent/60", className)}>
      <input type="checkbox" className="h-4 w-4 accent-primary" {...props} />
      <span>{label}</span>
    </label>
  );
}
