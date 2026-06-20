import type { ContainerSummary, ContainerDetail, ImageSummary, VolumeSummary, NetworkSummary } from "../types";

const TOKEN_KEY = "conmonitr_token";

function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

async function req<T>(url: string, init?: RequestInit): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    ...(init?.headers as Record<string, string>),
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(url, { ...init, headers });

  if (res.status === 401) {
    localStorage.removeItem(TOKEN_KEY);
    window.location.href = "/login";
    throw new Error("session expired");
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

// Build a same-origin WebSocket URL with the JWT attached as ?token=
export function wsUrl(path: string): string {
  const proto = window.location.protocol === "https:" ? "wss" : "ws";
  const token = getToken();
  const sep = path.includes("?") ? "&" : "?";
  const suffix = token ? `${sep}token=${token}` : "";
  return `${proto}://${window.location.host}${path}${suffix}`;
}
