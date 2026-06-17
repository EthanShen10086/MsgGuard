import { useEffect, useState } from "react";
import { fetchFlags, updateFlag } from "../api";

export default function Flags() {
  const [flags, setFlags] = useState<unknown[]>([]);
  const [error, setError] = useState("");
  const [key, setKey] = useState("cloud_llm");
  const [enabled, setEnabled] = useState(true);

  function load() {
    fetchFlags()
      .then(setFlags)
      .catch((e) => setError(e instanceof Error ? e.message : "failed"));
  }

  useEffect(() => {
    load();
  }, []);

  async function handleSave() {
    try {
      await updateFlag({ key, enabled, percentage: 100 });
      load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "save failed");
    }
  }

  return (
    <div>
      <h1>Feature Flags</h1>
      {error && <p className="error">{error}</p>}
      <div className="card">
        <input value={key} onChange={(e) => setKey(e.target.value)} placeholder="key" />
        <label>
          <input type="checkbox" checked={enabled} onChange={(e) => setEnabled(e.target.checked)} />
          enabled
        </label>
        <button type="button" onClick={handleSave}>Save</button>
      </div>
      <div className="card">
        <pre>{JSON.stringify(flags, null, 2)}</pre>
      </div>
    </div>
  );
}
