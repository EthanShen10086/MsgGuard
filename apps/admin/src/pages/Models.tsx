import { useEffect, useState } from "react";
import { fetchModelsLatest } from "../api";

export default function Models() {
  const [locale, setLocale] = useState("en");
  const [data, setData] = useState<Record<string, unknown> | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    setError("");
    fetchModelsLatest(locale)
      .then(setData)
      .catch((e) => setError(e instanceof Error ? e.message : "failed"));
  }, [locale]);

  return (
    <div>
      <h1>Models</h1>
      <label>
        Locale{" "}
        <input value={locale} onChange={(e) => setLocale(e.target.value)} />
      </label>
      {error && <p className="error">{error}</p>}
      {data && (
        <div className="card">
          <pre>{JSON.stringify(data, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}
