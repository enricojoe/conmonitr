import { useEffect, useState } from "react";
import { api } from "../api/client";
import type { NetworkSummary } from "../types";
import NavMenu from "../components/NavMenu";

export default function Networks() {
  const [networks, setNetworks] = useState<NetworkSummary[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api
      .listNetworks()
      .then(setNetworks)
      .catch((e: Error) => setError(e.message));
  }, []);

  return (
    <div className="mx-auto max-w-7xl px-4 py-6">
      <header className="mb-6">
        <h1 className="text-2xl font-bold text-zinc-100">ConMonitr</h1>
      </header>

      <NavMenu />

      <h2 className="mb-4 text-lg font-semibold text-zinc-200">
        Networks
        <span className="ml-2 text-sm font-normal text-zinc-500">
          {networks.length} total
        </span>
      </h2>

      {error && (
        <div className="mb-4 rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-2 text-sm text-rose-400">
          {error}
        </div>
      )}

      <div className="overflow-x-auto rounded-lg border border-zinc-800">
        <table className="w-full text-sm text-zinc-300">
          <thead className="border-b border-zinc-800 bg-zinc-900 text-xs uppercase tracking-wide text-zinc-500">
            <tr>
              <th className="px-4 py-3 text-left">Name</th>
              <th className="px-4 py-3 text-left">Driver</th>
              <th className="px-4 py-3 text-left">Scope</th>
              <th className="px-4 py-3 text-left">Internal</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-800">
            {networks.map((n) => (
              <tr
                key={n.id}
                className="bg-zinc-950 transition-colors hover:bg-zinc-900"
              >
                <td className="px-4 py-3 font-mono">{n.name}</td>
                <td className="px-4 py-3 text-zinc-400">{n.driver}</td>
                <td className="px-4 py-3 text-zinc-400">{n.scope}</td>
                <td className="px-4 py-3">
                  {n.internal ? (
                    <span className="inline-flex items-center rounded-full bg-amber-500/10 px-2 py-0.5 text-xs font-medium text-amber-400">
                      internal
                    </span>
                  ) : (
                    <span className="text-zinc-600">—</span>
                  )}
                </td>
              </tr>
            ))}
            {networks.length === 0 && !error && (
              <tr>
                <td
                  colSpan={4}
                  className="px-4 py-8 text-center text-zinc-600"
                >
                  No networks found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
