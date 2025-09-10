import {type ClassValue, clsx} from "clsx"
import {twMerge} from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const MouseButton = {
  Left: 0,
  Middle: 1,
  Right: 2,
} as const;
