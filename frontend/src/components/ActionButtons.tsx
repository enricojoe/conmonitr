import { useState } from "react";
import { api } from "../api/client";
import ConfirmModal from "./ConfirmModal";

interface Props {
  id: string;
  state: string;
  onChange: () => void;
}

type Action = "start" | "stop" | "restart" | "remove";

interface Pending {
  action: Action;
  fn: () => Promise<unknown>;
}

const MODAL: Record<
  Action,
  { title: string; message: string; confirmLabel: string; confirmClass: string }
> = {
  start: {
    title: "Start container?",
    message: "The container will be started.",
    confirmLabel: "Start",
    confirmClass: "bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30",
  },
  stop: {
    title: "Stop container?",
    message: "The container will be stopped. Running processes will be terminated.",
    confirmLabel: "Stop",
    confirmClass: "bg-amber-500/20 text-amber-400 hover:bg-amber-500/30",
  },
  restart: {
    title: "Restart container?",
    message: "The container will be restarted. Running processes will be briefly interrupted.",
    confirmLabel: "Restart",
    confirmClass: "bg-sky-500/20 text-sky-400 hover:bg-sky-500/30",
  },
  remove: {
    title: "Remove container?",
    message: "This will permanently remove the container and its writable layer. This cannot be undone.",
    confirmLabel: "Remove",
    confirmClass: "bg-rose-500/20 text-rose-400 hover:bg-rose-500/30",
  },
};

export default function ActionButtons({ id, state, onChange }: Props) {
  const [busy, setBusy] = useState<Action | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState<Pending | null>(null);
  const running = state === "running";

  const confirm = (action: Action, fn: () => Promise<unknown>) => {
    setPending({ action, fn });
  };

  const handleConfirm = async () => {
    if (!pending) return;
    const { action, fn } = pending;
    setPending(null);
    setBusy(action);
    setError(null);
    try {
      await fn();
      onChange();
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setBusy(null);
    }
  };

  const btn =
    "rounded-md px-2.5 py-1 text-xs font-medium transition-colors disabled:opacity-40 disabled:cursor-not-allowed";

  return (
    <>
      <div className="flex flex-wrap items-center gap-1.5">
        {running ? (
          <>
            <button
              className={`${btn} bg-amber-500/15 text-amber-400 hover:bg-amber-500/25`}
              disabled={busy !== null}
              onClick={() => confirm("stop", () => api.stop(id))}
            >
              {busy === "stop" ? "…" : "Stop"}
            </button>
            <button
              className={`${btn} bg-sky-500/15 text-sky-400 hover:bg-sky-500/25`}
              disabled={busy !== null}
              onClick={() => confirm("restart", () => api.restart(id))}
            >
              {busy === "restart" ? "…" : "Restart"}
            </button>
          </>
        ) : (
          <button
            className={`${btn} bg-emerald-500/15 text-emerald-400 hover:bg-emerald-500/25`}
            disabled={busy !== null}
            onClick={() => confirm("start", () => api.start(id))}
          >
            {busy === "start" ? "…" : "Start"}
          </button>
        )}
        <button
          className={`${btn} bg-rose-500/15 text-rose-400 hover:bg-rose-500/25`}
          disabled={busy !== null}
          onClick={() => confirm("remove", () => api.remove(id, true))}
        >
          {busy === "remove" ? "…" : "Remove"}
        </button>
        {error && <span className="text-xs text-rose-400">{error}</span>}
      </div>

      {pending && (
        <ConfirmModal
          {...MODAL[pending.action]}
          onConfirm={handleConfirm}
          onCancel={() => setPending(null)}
        />
      )}
    </>
  );
}
