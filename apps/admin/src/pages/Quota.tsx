import { useEffect, useState } from "react";
import { addQuotaWhitelist, fetchQuotaWhitelist } from "../api";

export default function Quota() {
  const [entries, setEntries] = useState<unknown[]>([]);
  const [error, setError] = useState("");
  const [userId, setUserId] = useState("");
  const [reason, setReason] = useState("");

  function load() {
    fetchQuotaWhitelist()
      .then(setEntries)
      .catch((e) => setError(e instanceof Error ? e.message : "failed"));
  }

  useEffect(() => {
    load();
  }, []);

  async function handleAdd() {
    if (!userId) return;
    try {
      await addQuotaWhitelist(userId, reason);
      setUserId("");
      setReason("");
      load();
    } catch (e) {
      setError(e instanceof Error ? e.message : "add failed");
    }
  }

  return (
    <div>
      <h1>Quota Whitelist</h1>
      {error && <p className="error">{error}</p>}
      <div className="card">
        <input value={userId} onChange={(e) => setUserId(e.target.value)} placeholder="user_id" />
        <input value={reason} onChange={(e) => setReason(e.target.value)} placeholder="reason" />
        <button type="button" onClick={handleAdd}>Add</button>
      </div>
      <div className="card">
        <pre>{JSON.stringify(entries, null, 2)}</pre>
      </div>
    </div>
  );
}
