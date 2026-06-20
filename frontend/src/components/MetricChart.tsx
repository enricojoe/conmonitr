import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import type { Metric } from "../types";

interface Props {
  title: string;
  data: Metric[];
  dataKey: keyof Metric;
  color: string;
  format?: (v: number) => string;
  unit?: string;
}

// MetricChart renders a single time-series area chart for one metric.
export default function MetricChart({
  title,
  data,
  dataKey,
  color,
  format,
  unit,
}: Props) {
  const points = data.map((m) => ({
    t: new Date(m.timestamp).toLocaleTimeString(),
    value: m[dataKey] as number,
  }));
  const current = points.length ? points[points.length - 1].value : 0;

  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900/50 p-4">
      <div className="mb-2 flex items-baseline justify-between">
        <h3 className="text-sm font-medium text-zinc-400">{title}</h3>
        <span className="text-lg font-semibold" style={{ color }}>
          {format ? format(current) : current.toFixed(2)}
          {unit && <span className="ml-1 text-xs text-zinc-500">{unit}</span>}
        </span>
      </div>
      <ResponsiveContainer width="100%" height={140}>
        <AreaChart data={points} margin={{ top: 4, right: 4, left: 4, bottom: 0 }}>
          <defs>
            <linearGradient id={`grad-${String(dataKey)}`} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor={color} stopOpacity={0.4} />
              <stop offset="100%" stopColor={color} stopOpacity={0} />
            </linearGradient>
          </defs>
          <XAxis dataKey="t" hide />
          <YAxis hide domain={[0, "auto"]} />
          <Tooltip
            contentStyle={{
              background: "#18181b",
              border: "1px solid #3f3f46",
              borderRadius: 8,
              fontSize: 12,
            }}
            formatter={(v: number) => [format ? format(v) : v.toFixed(2), title]}
          />
          <Area
            type="monotone"
            dataKey="value"
            stroke={color}
            strokeWidth={2}
            fill={`url(#grad-${String(dataKey)})`}
            isAnimationActive={false}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
