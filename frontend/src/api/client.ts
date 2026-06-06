// API client - centralized HTTP configuration
const API_BASE = '/api';

interface RequestOptions {
  method?: string;
  body?: any;
  params?: Record<string, string>;
  headers?: Record<string, string>;
}

export async function apiRequest<T = any>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, params, headers = {} } = options;

  let url = `${API_BASE}${path}`;
  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v !== '' && v !== undefined) searchParams.set(k, v);
    });
    const qs = searchParams.toString();
    if (qs) url += `?${qs}`;
  }

  const fetchOptions: RequestInit = {
    method,
    headers: {
      ...headers,
    },
  };

  if (body && method !== 'GET') {
    if (body instanceof FormData) {
      fetchOptions.body = body;
    } else {
      fetchOptions.headers = { ...fetchOptions.headers, 'Content-Type': 'application/json' };
      fetchOptions.body = JSON.stringify(body);
    }
  }

  const response = await fetch(url, fetchOptions);
  if (response.status === 204) return undefined as T;

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || `HTTP ${response.status}`);
  }

  return data as T;
}

// Helper to trigger browser file download
export function downloadFile(url: string, filename: string) {
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}
