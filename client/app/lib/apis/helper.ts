// Helper function to serialize query parameters
export function serializeQueryParams(params: Record<string, any>): string {
  if (!params || Object.keys(params).length === 0) return '';

  const searchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === null || value === undefined) return;

    if (Array.isArray(value)) {
      value.forEach((item) => {
        if (item !== null && item !== undefined) {
          searchParams.append(`${key}[]`, String(item));
        }
      });
    } else if (typeof value === 'object') {
      searchParams.append(key, JSON.stringify(value));
    } else {
      searchParams.append(key, String(value));
    }
  });

  return searchParams.toString();
}
