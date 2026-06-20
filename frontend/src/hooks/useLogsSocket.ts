import { useEffect, useState } from "react";
import { wsUrl, getToken } from "../api/client";
import type { LogLine } from "../types";

const MAX_LINES = 500;

// useLogsSocket subscribes to /ws/logs/{id} and keeps a bounded buffer of the
// most recent log lines.
export function useLogsSocket(id: string) {
  const [lines, setLines] = useState<LogLine[]>([]);

  useEffect(() => {
    setLines([]);
    let ws: WebSocket | null = null;
    let closed = false;
    let retry: ReturnType<typeof setTimeout>;

    const connect = () => {
      ws = new WebSocket(wsUrl(`/ws/logs/${id}?tail=200`));
      ws.onopen = () => ws!.send(getToken() ?? "");
      ws.onclose = () => {
        if (!closed) retry = setTimeout(connect, 2000);
      };
      ws.onmessage = (ev) => {
        const line: LogLine = JSON.parse(ev.data);
        setLines((prev) => [...prev, line].slice(-MAX_LINES));
      };
    };

    connect();
    return () => {
      closed = true;
      clearTimeout(retry);
      ws?.close();
    };
  }, [id]);

  return lines;
}
