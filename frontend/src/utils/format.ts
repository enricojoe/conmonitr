// Human-readable byte formatting.
export function formatBytes(bytes: number): string {
  if (!bytes || bytes < 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = bytes / Math.pow(1024, i);
  return `${value.toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

export function formatPercent(p: number): string {
  return `${p.toFixed(p < 10 ? 2 : 1)}%`;
}

export function shortId(id: string): string {
  return id.slice(0, 12);
}
