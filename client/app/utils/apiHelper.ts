import { GenericResponse } from '../lib/definitions';

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

export async function badRequestHandler(response: Response) {
  if (response.status > 299) {
    let errorMessage = 'An error occurred';
    try {
      const errorData = (await response.json()) as GenericResponse<unknown>;
      errorMessage =
        errorData.error?.details ||
        errorData.error?.stack ||
        errorData.error?.code ||
        'Unknown error';
    } catch (e) {
      const errorText = await response.text();
      errorMessage = errorText;
      console.error('Failed to parse error response:', e);
    }

    throw new Error(`Error ${response.status}: ${errorMessage}`);
  }
}
