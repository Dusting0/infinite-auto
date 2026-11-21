import type { ReactNode } from "react";
import { createPortal } from "react-dom";

export function DialogBackdrop({ children, onDismiss }: { children: ReactNode; onDismiss: () => void }) {
  // 脱离卡片的 backdrop-filter、overflow 与布局容器，确保弹窗覆盖视口且可点击。
  return createPortal(
    <div
      className="number-editor-backdrop"
      role="presentation"
      onPointerDown={(event) => {
        if (event.target === event.currentTarget) onDismiss();
      }}
    >
      {children}
    </div>,
    document.body,
  );
}
