import { useEffect, useRef, useState } from "react";
import { wsUrl, getToken } from "../api/client";
import type { Metric } from "../types";

const MAX_POINTS = 60;

// useAggregateStats subscribes to /ws/stats and maintains, per container, the
// latest metric plus a bounded rolling history for sparklines.
export function useAggregateStats() {
  const [latest, setLatest] = useState<Record<string, Metric>>({});
  const [history, setHistory] = useState<Record<string, Metric[]>>({});
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    let ws: WebSocket | null = null;
    let closed = false;
    let retry: ReturnType<typeof setTimeout>;

    const connect = () => {
      ws = new WebSocket(wsUrl("/ws/stats"));
      ws.onopen = () => {
        ws!.send(getToken() ?? "");
        setConnected(true);
      };
      ws.onclose = () => {
        setConnected(false);
        if (!closed) retry = setTimeout(connect, 2000);
      };
      ws.onmessage = (ev) => {
        const arr: Metric[] = JSON.parse(ev.data);
        const nextLatest: Record<string, Metric> = {};
        for (const m of arr) nextLatest[m.id] = m;
        setLatest(nextLatest);
        setHistory((prev) => {
          const next = { ...prev };
          for (const m of arr) {
            const series = next[m.id] ? [...next[m.id], m] : [m];
            next[m.id] = series.slice(-MAX_POINTS);
          }
          return next;
        });
      };
    };

    connect();
    return () => {
      closed = true;
      clearTimeout(retry);
      ws?.close();
    };
  }, []);

  return { latest, history, connected };
}

// useContainerStats subscribes to /ws/stats/{id} for one container and keeps a
// rolling time-series for the detail charts.
export function useContainerStats(id: string) {
  const [series, setSeries] = useState<Metric[]>([]);
  const idRef = useRef(id);
  idRef.current = id;

  useEffect(() => {
    setSeries([]);
    let ws: WebSocket | null = null;
    let closed = false;
    let retry: ReturnType<typeof setTimeout>;

    const connect = () => {
      ws = new WebSocket(wsUrl(`/ws/stats/${id}`));
      ws.onopen = () => ws!.send(getToken() ?? "");
      ws.onclose = () => {
        if (!closed) retry = setTimeout(connect, 2000);
      };
      ws.onmessage = (ev) => {
        const m: Metric = JSON.parse(ev.data);
        setSeries((prev) => [...prev, m].slice(-MAX_POINTS));
      };
    };

    connect();
    return () => {
      closed = true;
      clearTimeout(retry);
      ws?.close();
    };
  }, [id]);

  return series;
}
