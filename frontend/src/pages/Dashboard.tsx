import { useCallback, useEffect, useState } from "react";
import { api } from "../api/client";
import type { ContainerSummary } from "../types";
import { useAggregateStats } from "../hooks/useStatsSocket";
import ContainerGrid from "../components/ContainerGrid";
import NavMenu from "../components/NavMenu";

export default function Dashboard() {
  const [containers, setContainers] = useState<ContainerSummary[]>([]);
  const [error, setError] = useState<string | null>(null);
  const { latest, history, connected } = useAggregateStats();

  const refresh = useCallback(async () => {
    try {
      const list = await api.list();
      list.sort((a, b) => a.name.localeCompare(b.name));
      setContainers(list);
      setError(null);
    } catch (e) {
      setError((e as Error).message);
    }
  }, []);

  useEffect(() => {
    refresh();
    const t = setInterval(refresh, 5000);
    return () => clearInterval(t);
  }, [refresh]);

  const running = containers.filter((c) => c.state === "running").length;

  return (
    <div className="mx-auto max-w-7xl px-4 py-6">
      <header className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-zinc-100">ConMonitr</h1>
          <p className="text-sm text-zinc-500">
            {running} running · {containers.length} total
          </p>
        </div>
        <span
          className={`inline-flex items-center gap-1.5 text-xs ${
            connected ? "text-emerald-400" : "text-zinc-500"
          }`}
        >
          <span
            className={`h-2 w-2 rounded-full ${
              connected ? "bg-emerald-400" : "bg-zinc-600"
            }`}
          />
          {connected ? "live" : "connecting…"}
        </span>
      </header>

      <NavMenu />

      {error && (
        <div className="mb-4 rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-2 text-sm text-rose-400">
          {error}
        </div>
      )}

      <ContainerGrid
        containers={containers}
        latest={latest}
        history={history}
        onChange={refresh}
      />
    </div>
  );
}
