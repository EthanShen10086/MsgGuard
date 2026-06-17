import { useEffect, useState } from "react";
import { fetchFeedback } from "../api";

export default function Feedback() {
  const [items, setItems] = useState<unknown[]>([]);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchFeedback()
      .then(setItems)
      .catch((e) => setError(e instanceof Error ? e.message : "failed"));
  }, []);

  if (error) return <p className="error">{error}</p>;

  return (
    <div>
      <h1>Feedback</h1>
      <p>{items.length} items</p>
      <div className="card">
        <pre>{JSON.stringify(items.slice(0, 50), null, 2)}</pre>
      </div>
    </div>
  );
}
