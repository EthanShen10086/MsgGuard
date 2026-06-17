const GATEWAY = import.meta.env.VITE_GATEWAY_URL || "http://localhost:8080";
const TOKEN_KEY = "msgguard_admin_token";

export function getToken(): string {
  return sessionStorage.getItem(TOKEN_KEY) || "";
}

export function setToken(token: string): void {
  sessionStorage.setItem(TOKEN_KEY, token);
}

export function consumeTokenFromURL(): void {
  const params = new URLSearchParams(window.location.search);
  const t = params.get("access_token");
  if (!t) return;
  setToken(t);
  params.delete("access_token");
  const next = params.toString() ? `?${params}` : window.location.pathname;
  window.history.replaceState({}, "", next);
}

export async function fetchOIDCConfig(): Promise<{
  enabled: boolean;
  login_url: string;
  enforce_admin: boolean;
}> {
  return request("/api/v1/auth/oidc/config");
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  const token = getToken();
  if (token) headers.set("Authorization", `Bearer ${token}`);
  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json");
  }
  const res = await fetch(`${GATEWAY}${path}`, { ...init, headers });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${res.status}: ${text}`);
  }
  return res.json() as Promise<T>;
}

export async function fetchToken(userId: string): Promise<string> {
  const data = await request<{ access_token: string }>("/api/v1/auth/token", {
    method: "POST",
    body: JSON.stringify({ user_id: userId, roles: ["admin"] }),
  });
  return data.access_token;
}

export function fetchMetricsSummary() {
  return request<Record<string, unknown>>("/api/v1/admin/metrics/summary");
}

export function fetchFeedback() {
  return request<unknown[]>("/api/v1/feedback");
}

export function fetchModelsLatest(locale = "en") {
  return request<Record<string, unknown>>(`/api/v1/models/latest?locale=${locale}`);
}

export function fetchFlags() {
  return request<unknown[]>("/api/v1/admin/flags");
}

export function updateFlag(flag: Record<string, unknown>) {
  return request<{ status: string }>("/api/v1/admin/flags", {
    method: "POST",
    body: JSON.stringify(flag),
  });
}

export function fetchQuotaWhitelist() {
  return request<unknown[]>("/api/v1/admin/quota/whitelist");
}

export function addQuotaWhitelist(userId: string, reason: string) {
  return request<{ status: string }>("/api/v1/admin/quota/whitelist", {
    method: "POST",
    body: JSON.stringify({ user_id: userId, reason }),
  });
}
