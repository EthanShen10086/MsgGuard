import { useEffect, useState } from "react";
import { fetchMetricsSummary } from "../api";

export default function Dashboard() {
  const [data, setData] = useState<Record<string, unknown> | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchMetricsSummary()
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : "failed"));
  }, []);

  if (error) return <p className="error">{error}</p>;
  if (!data) return <p>Loading…</p>;

  const counts = (data.event_counts as Record<string, number>) || {};

  return (
    <div>
      <h1>Dashboard</h1>
      <div className="card">
        <p>Period: {String(data.period_days)} days</p>
        <p>Feedback total: {String(data.feedback_total)}</p>
        <p>Shadow total: {String(data.shadow_total)}</p>
        <p>Shadow disagree rate: {String(data.shadow_disagree_rate)}</p>
      </div>
      <div className="card">
        <h3>Event counts</h3>
        <ul>
          {Object.entries(counts).map(([k, v]) => (
            <li key={k}>{k}: {v}</li>
          ))}
          {Object.keys(counts).length === 0 && <li>No events</li>}
        </ul>
      </div>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}
