import type { ContainerSummary, ContainerDetail, ImageSummary, VolumeSummary, NetworkSummary } from "../types";

async function req<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
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

// Build a same-origin WebSocket URL (works through the Vite dev proxy).
export function wsUrl(path: string): string {
  const proto = window.location.protocol === "https:" ? "wss" : "ws";
  return `${proto}://${window.location.host}${path}`;
}
