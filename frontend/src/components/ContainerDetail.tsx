import type { ContainerDetail } from "../types";
import { shortId } from "../utils/format";

interface Props {
  detail: ContainerDetail;
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-4">
      <h3 className="mb-2 text-sm font-medium text-zinc-400">{title}</h3>
      {children}
    </div>
  );
}

function Row({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex justify-between gap-4 py-0.5 text-sm">
      <span className="text-zinc-500">{label}</span>
      <span className="truncate text-right font-mono text-zinc-300">{value}</span>
    </div>
  );
}

export default function ContainerDetailPanel({ detail }: Props) {
  const ports = detail.ports ?? [];
  const networks = detail.networks ?? [];
  const mounts = detail.mounts ?? [];
  const env = detail.env ?? [];

  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <Section title="Overview">
        <Row label="ID" value={shortId(detail.id)} />
        <Row label="Image" value={detail.image} />
        <Row label="State" value={detail.state} />
        <Row label="Restart policy" value={detail.restartPolicy || "no"} />
        <Row label="TTY" value={detail.tty ? "yes" : "no"} />
        <Row
          label="Command"
          value={(detail.command ?? []).join(" ") || "—"}
        />
      </Section>

      <Section title="Networks">
        {networks.length === 0 ? (
          <p className="text-sm text-zinc-600">none</p>
        ) : (
          networks.map((n) => (
            <Row key={n.name} label={n.name} value={n.ipAddress || "—"} />
          ))
        )}
      </Section>

      <Section title="Ports">
        {ports.length === 0 ? (
          <p className="text-sm text-zinc-600">none</p>
        ) : (
          ports.map((p, i) => (
            <Row
              key={i}
              label={`${p.privatePort}/${p.type}`}
              value={p.publicPort ? `→ ${p.ip || "0.0.0.0"}:${p.publicPort}` : "internal"}
            />
          ))
        )}
      </Section>

      <Section title="Mounts">
        {mounts.length === 0 ? (
          <p className="text-sm text-zinc-600">none</p>
        ) : (
          mounts.map((m, i) => (
            <Row
              key={i}
              label={m.destination}
              value={`${m.type}${m.rw ? "" : " (ro)"}`}
            />
          ))
        )}
      </Section>

      <div className="lg:col-span-2">
        <Section title={`Environment (${env.length})`}>
          {env.length === 0 ? (
            <p className="text-sm text-zinc-600">none</p>
          ) : (
            <div className="max-h-48 overflow-auto font-mono text-xs text-zinc-400">
              {env.map((e, i) => (
                <div key={i} className="truncate py-0.5">
                  {e}
                </div>
              ))}
            </div>
          )}
        </Section>
      </div>
    </div>
  );
}
