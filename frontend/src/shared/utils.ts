/**
 * Shared utility functions
 * Can be used by main, preload, and renderer processes
 */

export const formatDate = (date: string | Date): string => {
  return new Date(date).toLocaleDateString();
};

export const formatNumber = (num: number): string => {
  return new Intl.NumberFormat().format(num);
};

export const sleep = (ms: number): Promise<void> => new Promise(resolve => setTimeout(resolve, ms));

