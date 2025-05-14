export const capitalize = (s: string) =>
  (s && s[0].toUpperCase() + s.slice(1)) || '';

/**
 * Format number as currency with 2 decimal places
 * @param value Number to format
 * @returns Formatted string (e.g. 1,234.56)
 */
export const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
};

/**
 * Format date in a human-readable format
 * @param date Date to format
 * @returns Formatted date string (e.g. May 16, 2023)
 */
export const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    month: 'long',
    day: 'numeric',
    year: 'numeric'
  }).format(date);
};
