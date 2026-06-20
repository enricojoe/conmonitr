import { Link } from "react-router-dom";
import type { ContainerSummary, Metric } from "../types";
import { formatBytes, formatPercent, shortId } from "../utils/format";
import StatusBadge from "./StatusBadge";
import ActionButtons from "./ActionButtons";
import Sparkline from "./Sparkline";

interface Props {
  containers: ContainerSummary[];
  latest: Record<string, Metric>;
  history: Record<string, Metric[]>;
  onChange: () => void;
}

export default function ContainerGrid({
  containers,
  latest,
  history,
  onChange,
}: Props) {
  if (containers.length === 0) {
    return (
      <p className="py-16 text-center text-zinc-500">No containers found.</p>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
      {containers.map((c) => {
        const m = latest[c.id];
        const cpuHist = (history[c.id] ?? []).map((p) => p.cpuPercent);
        const memHist = (history[c.id] ?? []).map((p) => p.memUsage);
        return (
          <div
            key={c.id}
            className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-4 transition-colors hover:border-zinc-700"
          >
            <div className="flex items-start justify-between gap-2">
              <div className="min-w-0">
                <Link
                  to={`/containers/${c.id}`}
                  className="block truncate font-semibold text-zinc-100 hover:text-sky-400"
                >
                  {c.name}
                </Link>
                <p className="truncate text-xs text-zinc-500">{c.image}</p>
                <p className="font-mono text-[10px] text-zinc-600">
                  {shortId(c.id)}
                </p>
              </div>
              <StatusBadge state={c.state} />
            </div>

            <div className="mt-3 grid grid-cols-2 gap-3">
              <div>
                <div className="flex items-baseline justify-between">
                  <span className="text-xs text-zinc-500">CPU</span>
                  <span className="text-xs font-medium text-sky-400">
                    {m ? formatPercent(m.cpuPercent) : "—"}
                  </span>
                </div>
                <Sparkline values={cpuHist} color="#38bdf8" />
              </div>
              <div>
                <div className="flex items-baseline justify-between">
                  <span className="text-xs text-zinc-500">Memory</span>
                  <span className="text-xs font-medium text-violet-400">
                    {m ? formatBytes(m.memUsage) : "—"}
                  </span>
                </div>
                <Sparkline values={memHist} color="#a78bfa" />
              </div>
            </div>

            <div className="mt-3 flex items-center justify-between border-t border-zinc-800 pt-3">
              <span className="truncate text-xs text-zinc-500">{c.status}</span>
              <ActionButtons id={c.id} state={c.state} onChange={onChange} />
            </div>
          </div>
        );
      })}
    </div>
  );
}
