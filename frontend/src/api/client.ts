import type { ContainerSummary, ContainerDetail, ImageSummary, VolumeSummary, NetworkSummary } from "../types";

const TOKEN_KEY = "conmonitr_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

// Deduplicate concurrent refresh calls — only one in-flight at a time.
let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  if (refreshPromise) return refreshPromise;

  refreshPromise = fetch("/api/auth/refresh", { method: "POST" })
    .then(async (res) => {
      if (!res.ok) throw new Error("refresh failed");
      const data = await res.json();
      setToken(data.token);
      return data.token as string;
    })
    .finally(() => {
      refreshPromise = null;
    });

  return refreshPromise;
}

async function req<T>(url: string, init?: RequestInit): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(init?.headers as Record<string, string>),
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(url, { ...init, headers });

  if (res.status === 401) {
    try {
      const newToken = await refreshAccessToken();
      const retryHeaders = { ...headers, Authorization: `Bearer ${newToken}` };
      const retry = await fetch(url, { ...init, headers: retryHeaders });
      if (!retry.ok) {
        const body = await retry.json().catch(() => ({}));
        throw new Error(body.error || `request failed: ${retry.status}`);
      }
      return retry.json() as Promise<T>;
    } catch {
      clearToken();
      window.location.href = "/login";
      throw new Error("session expired");
    }
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `request failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  list: () => req<ContainerSummary[]>("/api/containers/"),
  inspect: (id: string) => req<ContainerDetail>(`/api/containers/${id}`),
  start: (id: string) =>
    req(`/api/containers/${id}/start`, { method: "POST" }),
  stop: (id: string) => req(`/api/containers/${id}/stop`, { method: "POST" }),
  restart: (id: string) =>
    req(`/api/containers/${id}/restart`, { method: "POST" }),
  remove: (id: string, force = false) =>
    req(`/api/containers/${id}?force=${force}`, { method: "DELETE" }),
  listImages: () => req<ImageSummary[]>("/api/images"),
  listVolumes: () => req<VolumeSummary[]>("/api/volumes"),
  listNetworks: () => req<NetworkSummary[]>("/api/networks"),
};

// Build a same-origin WebSocket URL. Auth is done via the first message, not the URL.
export function wsUrl(path: string): string {
  const proto = window.location.protocol === "https:" ? "wss" : "ws";
  return `${proto}://${window.location.host}${path}`;
}
