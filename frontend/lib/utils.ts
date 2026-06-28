import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

// cn は条件付きクラス結合と Tailwind の競合解決をまとめて行う（shadcn/ui 標準ヘルパー）。
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
