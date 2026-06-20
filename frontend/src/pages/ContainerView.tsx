import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api/client";
import type { ContainerDetail } from "../types";
import { useContainerStats } from "../hooks/useStatsSocket";
import { useLogsSocket } from "../hooks/useLogsSocket";
import { formatBytes, formatPercent } from "../utils/format";
import MetricChart from "../components/MetricChart";
import ContainerDetailPanel from "../components/ContainerDetail";
import LogViewer from "../components/LogViewer";
import StatusBadge from "../components/StatusBadge";
import ActionButtons from "../components/ActionButtons";

export default function ContainerView() {
  const { id = "" } = useParams();
  const [detail, setDetail] = useState<ContainerDetail | null>(null);
  const [error, setError] = useState<string | null>(null);
  const series = useContainerStats(id);
  const logs = useLogsSocket(id);

  const refresh = useCallback(async () => {
    try {
      setDetail(await api.inspect(id));
      setError(null);
    } catch (e) {
      setError((e as Error).message);
    }
  }, [id]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return (
    <div className="mx-auto max-w-7xl px-4 py-6">
      <Link to="/" className="text-sm text-sky-400 hover:text-sky-300">
        ← Back to dashboard
      </Link>

      <header className="mt-3 mb-6 flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-zinc-100">
            {detail?.name ?? id.slice(0, 12)}
          </h1>
          {detail && <StatusBadge state={detail.state} />}
        </div>
        {detail && (
          <ActionButtons id={id} state={detail.state} onChange={refresh} />
        )}
      </header>

      {error && (
        <div className="mb-4 rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-2 text-sm text-rose-400">
          {error}
        </div>
      )}

      <div className="mb-6 grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricChart
          title="CPU"
          data={series}
          dataKey="cpuPercent"
          color="#38bdf8"
          format={formatPercent}
        />
        <MetricChart
          title="Memory"
          data={series}
          dataKey="memUsage"
          color="#a78bfa"
          format={formatBytes}
        />
        <MetricChart
          title="Network RX"
          data={series}
          dataKey="netRx"
          color="#34d399"
          format={formatBytes}
        />
        <MetricChart
          title="Block Read"
          data={series}
          dataKey="blkRead"
          color="#fbbf24"
          format={formatBytes}
        />
      </div>

      {detail && (
        <div className="mb-6">
          <ContainerDetailPanel detail={detail} />
        </div>
      )}

      <h2 className="mb-2 text-sm font-medium text-zinc-400">Logs</h2>
      <LogViewer lines={logs} />
    </div>
  );
}
