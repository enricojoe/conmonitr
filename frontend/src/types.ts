// Shared types mirroring the Go backend JSON DTOs.

export interface Port {
  ip?: string;
  privatePort: number;
  publicPort?: number;
  type: string;
}

export interface ContainerSummary {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  created: number;
  ports: Port[] | null;
}

export interface NetworkInfo {
  name: string;
  ipAddress: string;
  gateway: string;
  macAddress: string;
}

export interface MountInfo {
  type: string;
  source: string;
  destination: string;
  name?: string;
  rw: boolean;
}

export interface ContainerDetail {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  created: string;
  command: string[] | null;
  tty: boolean;
  restartPolicy: string;
  env: string[] | null;
  ports: Port[] | null;
  networks: NetworkInfo[] | null;
  mounts: MountInfo[] | null;
}

export interface Metric {
  id: string;
  name: string;
  timestamp: number;
  cpuPercent: number;
  memUsage: number;
  memLimit: number;
  memPercent: number;
  netRx: number;
  netTx: number;
  blkRead: number;
  blkWrite: number;
}

export interface LogLine {
  stream: "stdout" | "stderr";
  line: string;
  timestamp: number;
}

export interface ImageSummary {
  id: string;
  repoTags: string[] | null;
  size: number;
  created: number;
  containers: number;
}

export interface VolumeSummary {
  name: string;
  driver: string;
  mountpoint: string;
  scope: string;
  createdAt: string;
}

export interface NetworkSummary {
  id: string;
  name: string;
  driver: string;
  scope: string;
  created: string;
  internal: boolean;
}
