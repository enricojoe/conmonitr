interface Props {
  state: string;
}

const COLORS: Record<string, string> = {
  running: "bg-emerald-500/15 text-emerald-400 ring-emerald-500/30",
  exited: "bg-rose-500/15 text-rose-400 ring-rose-500/30",
  paused: "bg-amber-500/15 text-amber-400 ring-amber-500/30",
  created: "bg-sky-500/15 text-sky-400 ring-sky-500/30",
  restarting: "bg-amber-500/15 text-amber-400 ring-amber-500/30",
};

export default function StatusBadge({ state }: Props) {
  const cls = COLORS[state] ?? "bg-zinc-500/15 text-zinc-400 ring-zinc-500/30";
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset ${cls}`}
    >
      <span className="h-1.5 w-1.5 rounded-full bg-current" />
      {state}
    </span>
  );
}
