import { useEffect, useState } from "react";
import { NavLink, Route, Routes } from "react-router-dom";
import { fetchToken, getToken, setToken } from "./api";
import Dashboard from "./pages/Dashboard";
import Feedback from "./pages/Feedback";
import Models from "./pages/Models";
import Flags from "./pages/Flags";
import Quota from "./pages/Quota";

export default function App() {
  const [token, setTokenState] = useState(getToken());
  const [userId, setUserId] = useState("admin");

  useEffect(() => {
    setTokenState(getToken());
  }, []);

  async function handleLogin() {
    try {
      const t = await fetchToken(userId);
      setToken(t);
      setTokenState(t);
    } catch (e) {
      alert(e instanceof Error ? e.message : "login failed");
    }
  }

  return (
    <>
      <div className="auth-bar">
        <strong>MsgGuard Admin</strong>
        <input
          value={userId}
          onChange={(e) => setUserId(e.target.value)}
          placeholder="user_id"
        />
        <button type="button" onClick={handleLogin}>
          Get token (dev)
        </button>
        <span className={token ? "" : "error"}>
          {token ? "Authenticated" : "No token — admin API calls will 401"}
        </span>
      </div>
      <div className="layout">
        <nav>
          <h2>Admin</h2>
          <NavLink to="/" end>Dashboard</NavLink>
          <NavLink to="/feedback">Feedback</NavLink>
          <NavLink to="/models">Models</NavLink>
          <NavLink to="/flags">Flags</NavLink>
          <NavLink to="/quota">Quota</NavLink>
        </nav>
        <main>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/feedback" element={<Feedback />} />
            <Route path="/models" element={<Models />} />
            <Route path="/flags" element={<Flags />} />
            <Route path="/quota" element={<Quota />} />
          </Routes>
        </main>
      </div>
    </>
  );
}
