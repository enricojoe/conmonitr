import { Line, LineChart, ResponsiveContainer, YAxis } from "recharts";

interface Props {
  values: number[];
  color: string;
}

// Sparkline renders a tiny inline trend line for a metric history.
export default function Sparkline({ values, color }: Props) {
  const data = values.map((v, i) => ({ i, v }));
  return (
    <ResponsiveContainer width="100%" height={32}>
      <LineChart data={data}>
        <YAxis hide domain={[0, "auto"]} />
        <Line
          type="monotone"
          dataKey="v"
          stroke={color}
          strokeWidth={1.5}
          dot={false}
          isAnimationActive={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
