import { useState, useEffect, createContext, useContext } from "react";


const BASE = "http://localhost:8080";

async function apiFetch(path, opts = {}) {
  const token = localStorage.getItem("token");

  const headers = {
    "Content-Type": "application/json",
    ...opts.headers,
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const res = await fetch(BASE + path, {
    ...opts,
    headers,
  });

  if (res.status === 401) {
    localStorage.removeItem("token");
    localStorage.removeItem("userId");

    window.dispatchEvent(new Event("unauthorized"));

    throw new Error("Session expired");
  }

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `HTTP ${res.status}`);
  }

  if (
    res.status === 204 ||
    res.headers.get("content-length") === "0"
  ) {
    return null;
  }

  const text = await res.text();

  if (!text) {
    return null;
  }

  return JSON.parse(text);
}

const api = {
  register: (email, password) =>
    apiFetch("/auth/register", { method: "POST", body: JSON.stringify({ email, password }) }),
  login: (email, password) =>
    apiFetch("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),
  logout: () => apiFetch("/auth/logout", { method: "POST" }),
  vehicles: () => apiFetch("/vehicles"),
  addVehicle: (plate_number) =>
    apiFetch("/vehicles", { method: "POST", body: JSON.stringify({ plate_number }) }),
  deleteVehicle: (plate_number) =>
    apiFetch("/vehicles", {
      method: "DELETE",
      body: JSON.stringify({ plate_number }),
    }),
  available: () => apiFetch("/slots/available"),
  startSession: (vehicle_id) =>
    apiFetch("/sessions/start", { method: "POST", body: JSON.stringify({ vehicle_id }) }),
  endSession: (id) => apiFetch(`/sessions/${id}/end`, { method: "POST" }),
  getSession: (id) => apiFetch(`/sessions/${id}`),
  getPrice: (id) => apiFetch(`/sessions/${id}/price`),
  activeSessions: () =>
    apiFetch("/sessions/active"),
  // historySessions: () => apiFetch("/sessions/history"),

  adminSlots: () => apiFetch("/admin/slots"),

  addSlots: (count) =>
    apiFetch("/admin/slots", {
      method: "POST",
      body: JSON.stringify({ count }),
    }),
};

const AuthCtx = createContext(null);
function useAuth() { return useContext(AuthCtx); }


const Icon = ({ d, size = 20, stroke = "currentColor" }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke={stroke} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
    <path d={d} />
  </svg>
);

const Icons = {
  car:     "M5 17H3a2 2 0 01-2-2V9a2 2 0 012-2h14a2 2 0 012 2v6a2 2 0 01-2 2h-2M5 17a2 2 0 104 0 2 2 0 00-4 0zm10 0a2 2 0 104 0 2 2 0 00-4 0zM3 9l1-4h12l1 4",
  parking: "M9 12h6m-6 4h6M5 3a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2H5zm4 4h3a2 2 0 010 4h-3V7z",
  logout:  "M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1",
  plus:    "M12 4v16m-8-8h16",
  trash:   "M3 6h18M8 6V4h8v2M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6",
  clock:   "M12 2a10 10 0 110 20A10 10 0 0112 2zm0 4v6l4 2",
  check:   "M20 6L9 17l-5-5",
  user:    "M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 3a4 4 0 110 8 4 4 0 010-8z",
  bolt:    "M13 2L3 14h9l-1 8 10-12h-9l1-8z",
};

function fmtTime(unix) {
  if (!unix) return "—";
  return new Date(unix * 1000).toLocaleString("ru-RU", {
    day: "2-digit", month: "2-digit", hour: "2-digit", minute: "2-digit",
  });
}

function fmtDuration(seconds) {
  if (!seconds) return "—";
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  return h > 0 ? `${h}ч ${m}м` : `${m}м`;
}

function Spinner() {
  return (
    <div style={{ display: "flex", justifyContent: "center", padding: "2rem" }}>
      <div className="spinner" />
    </div>
  );
}


function Toast({ msg, type, onClose }) {
  useEffect(() => {
    const t = setTimeout(onClose, 3500);
    return () => clearTimeout(t);
  }, [msg]);
  if (!msg) return null;
  return (
    <div className={`toast toast-${type}`}>
      {type === "ok" ? <Icon d={Icons.check} size={16} /> : "⚠"}
      <span>{msg}</span>
    </div>
  );
}

function Badge({ children, color = "blue" }) {
  return <span className={`badge badge-${color}`}>{children}</span>;
}

function AuthPage({ onAuth }) {
  const [mode, setMode] = useState("login");
  const [email, setEmail] = useState("");
  const [pass, setPass] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);


  async function submit(e) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const fn = mode === "login" ? api.login : api.register;
      const data = await fn(email, pass);
      if (!data.token) {
        throw new Error("Token not received");
      }

      localStorage.setItem("token", data.token);
      localStorage.setItem("userId", data.user_id);
      localStorage.setItem("userEmail", email);
      onAuth(data, email);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-wrap">
      <div className="auth-card">
        <div className="auth-logo">
          <div className="logo-icon"><Icon d={Icons.parking} size={28} stroke="#fff" /></div>
          <div>
            <div className="logo-title">SmartPark</div>
            <div className="logo-sub">система управления парковкой</div>
          </div>
        </div>

        <div className="tab-group">
          <button className={`tab ${mode === "login" ? "active" : ""}`} onClick={() => setMode("login")}>Войти</button>
          <button className={`tab ${mode === "register" ? "active" : ""}`} onClick={() => setMode("register")}>Регистрация</button>
        </div>

        <form onSubmit={submit} className="auth-form">
          <label className="field-label">Email</label>
          <input className="field-input" type="email" value={email} onChange={e => setEmail(e.target.value)} placeholder="you@example.com" required />

          <label className="field-label" style={{ marginTop: "1rem" }}>Пароль</label>
          <input className="field-input" type="password" value={pass} onChange={e => setPass(e.target.value)} placeholder="••••••••" required minLength={6} />

          {error && <div className="auth-error">{error}</div>}

          <button className="btn-primary" type="submit" disabled={loading} style={{ marginTop: "1.5rem", width: "100%" }}>
            {loading ? "..." : mode === "login" ? "Войти" : "Создать аккаунт"}
          </button>
        </form>
      </div>
    </div>
  );
}

function SlotsWidget() {
  const [count, setCount] = useState(null);
  const [loading, setLoading] = useState(true);

  async function load() {
    try {
      const d = await api.available();
      setCount(d.count);
    } catch (e) {
      if (e.message !== "Session expired") {
        console.error(e);
      }
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    const t = setInterval(load, 10000);
    return () => clearInterval(t);
  }, []);

  const pct = count != null ? Math.round((count / 50) * 100) : 0;
  const color = count > 15 ? "#22c55e" : count > 5 ? "#f59e0b" : "#ef4444";

  return (
    <div className="slots-widget">
      <div className="slots-header">
        <Icon d={Icons.bolt} size={18} />
        <span>Свободных мест</span>
      </div>
      {loading ? <Spinner /> : (
        <>
          <div className="slots-number" style={{ color }}>{count ?? "—"}</div>
          <div className="slots-sub">из 50 мест</div>
          <div className="slots-bar-wrap">
            <div className="slots-bar" style={{ width: `${pct}%`, background: color }} />
          </div>
          <button className="btn-ghost" onClick={load} style={{ marginTop: "0.5rem", fontSize: "12px" }}>↻ обновить</button>
        </>
      )}
    </div>
  );
}
function VehiclesPage({ onStartSession, toast }) {

  const [vehicles, setVehicles] = useState([])
  const [plate, setPlate] = useState("")
  const [loading, setLoading] = useState(true)
  const [adding, setAdding] = useState(false)

  async function load() {

    try {

      const d = await api.vehicles()

      setVehicles(d.vehicles || [])

    } catch (e) {

      toast(e.message, "err")

    } finally {

      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  async function addVehicle(e) {

    e.preventDefault()

    if (!plate.trim()) return

    setAdding(true)

    try {

      await api.addVehicle(
        plate.trim().toUpperCase()
      )

      setPlate("")

      toast("Автомобиль добавлен", "ok")

      load()

    } catch (e) {

      toast(e.message, "err")

    } finally {

      setAdding(false)
    }
  }

  async function deleteVehicle(plate) {

    try {

      await api.deleteVehicle(plate)

      toast("Автомобиль удален", "ok")

      load()

    } catch (e) {

      toast(e.message, "err")
    }
  }

  return (
    <div className="page">

      <h2 className="page-title">
        Мои автомобили
      </h2>

      <form
        onSubmit={addVehicle}
        className="add-form"
      >

        <input
          className="field-input"
          value={plate}
          onChange={(e) => setPlate(e.target.value)}
          placeholder="ABC-123"
          style={{ flex: 1 }}
        />

        <button
          className="btn-primary"
          type="submit"
          disabled={adding}
        >
          <Icon d={Icons.plus} size={16} />
          Добавить
        </button>

      </form>

      {loading ? (

        <Spinner />

      ) : vehicles.length === 0 ? (

        <div className="empty-state">
          <Icon
            d={Icons.car}
            size={40}
            stroke="#94a3b8"
          />
          <p>Нет автомобилей</p>
        </div>

      ) : (

        <div className="card-grid">

          {vehicles.map((v) => (

            <div
              className="vehicle-card"
              key={v}
            >

              <div className="vehicle-plate">
                {v}
              </div>

              <div className="vehicle-actions">

                <button
                  className="btn-success"
                  onClick={() => onStartSession(v)}
                >
                  <Icon
                    d={Icons.parking}
                    size={15}
                  />

                  Припарковать
                </button>

                <button
                  className="btn-danger-ghost"
                  onClick={() => deleteVehicle(v)}
                >
                  <Icon
                    d={Icons.trash}
                    size={15}
                  />
                </button>

              </div>

            </div>

          ))}

        </div>

      )}

    </div>
  )
}
function StartSessionModal({
  vehicleId,
  vehicles,
  onClose,
  onStarted,
  toast,
}) {

  const [selectedVehicle, setSelectedVehicle] = useState(
    vehicleId || null
  )

  const [loading, setLoading] = useState(false)

  async function start() {

    if (!selectedVehicle) {

      toast("Выберите автомобиль", "err")

      return
    }

    setLoading(true)

    try {

      // backend currently returns []string

      const idx = vehicles.indexOf(selectedVehicle)

      const vehicleIdNum =
        idx >= 0 ? idx + 1 : 1

      const sess = await api.startSession(
        vehicleIdNum
      )

      onStarted(sess)

      toast(
        "Парковка начата! Место #" +
          sess.slot_id,
        "ok"
      )

      onClose()

    } catch (e) {

      toast(e.message, "err")

    } finally {

      setLoading(false)
    }
  }

  return (
    <div
      className="modal-backdrop"
      onClick={onClose}
    >

      <div
        className="modal"
        onClick={(e) => e.stopPropagation()}
      >

        <h3 className="modal-title">
          Начать парковку
        </h3>

        <p className="modal-sub">
          Выберите автомобиль:
        </p>

        <div className="modal-options">

          {vehicles.map((v) => (

            <button
              key={v}
              className={`modal-option ${
                selectedVehicle === v
                  ? "selected"
                  : ""
              }`}
              onClick={() =>
                setSelectedVehicle(v)
              }
            >

              <Icon
                d={Icons.car}
                size={16}
              />

              {v}

            </button>

          ))}

        </div>

        <div className="modal-actions">

          <button
            className="btn-ghost"
            onClick={onClose}
          >
            Отмена
          </button>

          <button
            className="btn-primary"
            onClick={start}
            disabled={loading || !selectedVehicle}
          >
            {loading ? "..." : "Начать"}
          </button>

        </div>

      </div>

    </div>
  )
}

function SessionsPage({ activeSessions, onEnd, toast }) {
  const [prices, setPrices] = useState({});

  async function loadPrice(id) {
    try {
      const d = await api.getPrice(id);
      setPrices(p => ({ ...p, [id]: d }));
    } catch {}
  }

  useEffect(() => {
    activeSessions.forEach(s => loadPrice(s.session_id));
    const t = setInterval(() => activeSessions.forEach(s => loadPrice(s.session_id)), 30000);
    return () => clearInterval(t);
  }, [activeSessions]);

  async function endSession(id) {
    try {
      const d = await api.endSession(id);
      toast(`Сессия завершена. Сумма: ${d.amount} ₸`, "ok");
      onEnd(id, d);
    } catch (e) { toast(e.message, "err"); }
  }

  if (activeSessions.length === 0) {
    return (
      <div className="page">
        <h2 className="page-title">Активные парковки</h2>
        <div className="empty-state">
          <Icon d={Icons.clock} size={40} stroke="#94a3b8" />
          <p>Нет активных парковочных сессий</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <h2 className="page-title">Активные парковки</h2>
      <div className="session-list">
        {activeSessions.map(s => {
          const p = prices[s.session_id];
          return (
            <div className="session-card" key={s.session_id}>
              <div className="session-top">
                <div>
                  <div className="session-id">Сессия #{s.session_id}</div>
                  <div className="session-meta">Место <strong>#{s.slot_id}</strong> · Начало {fmtTime(s.start_time_unix)}</div>
                </div>
                <Badge color="green">Активна</Badge>
              </div>
              {p && (
                <div className="session-price-row">
                  <div className="session-price-item">
                    <span className="session-price-label">Время</span>
                    <span className="session-price-val">{fmtDuration(p.elapsed_seconds)}</span>
                  </div>
                  <div className="session-price-item">
                    <span className="session-price-label">Текущая сумма</span>
                    <span className="session-price-val accent">{p.amount} ₸</span>
                  </div>
                </div>
              )}
              <button className="btn-end" onClick={() => endSession(s.session_id)}>
                <Icon d={Icons.check} size={15} /> Завершить и оплатить
              </button>
            </div>
          );
        })}
      </div>
    </div>
  );
}


function HistoryPage({ history }) {
  if (history.length === 0) {
    return (
      <div className="page">
        <h2 className="page-title">История парковок</h2>
        <div className="empty-state">
          <Icon d={Icons.clock} size={40} stroke="#94a3b8" />
          <p>История пуста</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <h2 className="page-title">История парковок</h2>
      <div className="session-list">
        {[...history].reverse().map(s => (
          <div className="session-card completed" key={s.session_id + "-hist"}>
            <div className="session-top">
              <div>
                <div className="session-id">Сессия #{s.session_id}</div>
                <div className="session-meta">Место #{s.slot_id} · {fmtTime(s.start_time_unix)} — {fmtTime(s.end_time_unix)}</div>
              </div>
              <Badge color="gray">Завершена</Badge>
            </div>
            {s.amount > 0 && (
              <div className="session-price-row">
                <div className="session-price-item">
                  <span className="session-price-label">Итого оплачено</span>
                  <span className="session-price-val accent">{s.amount} ₸</span>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

function ProfilePage({ auth, vehicles, activeSessions, history }) {
  return (
    <div className="page">
      <h2 className="page-title">Личный кабинет</h2>

      <div className="session-card">
        <div className="session-id">Профиль</div>
        <div className="session-meta">Email: {auth.email || "—"}</div>
        <div className="session-meta">User ID: {auth.user_id}</div>
      </div>

      <div className="session-card">
        <div className="session-id">Статистика</div>
        <div className="session-meta">Автомобилей: {vehicles.length}</div>
        <div className="session-meta">Активных парковок: {activeSessions.length}</div>
        <div className="session-meta">Завершённых парковок: {history.length}</div>
      </div>
    </div>
  );
}

function AdminSlotsPage({ toast }) {
  const [slots, setSlots] = useState([]);
  const [count, setCount] = useState(1);
  const [loading, setLoading] = useState(true);

  async function load() {
    try {
      const d = await api.adminSlots();
      setSlots(d.slots || []);
    } catch (e) {
      toast(e.message, "err");
    } finally {
      setLoading(false);
    }
  }

  async function add(e) {
    e.preventDefault();

    try {
      await api.addSlots(Number(count));
      toast("Места добавлены", "ok");
      setCount(1);
      load();
    } catch (e) {
      toast(e.message, "err");
    }
  }

  useEffect(() => {
    load();
  }, []);

  return (
    <div className="page">
      <h2 className="page-title">Админ: парковочные места</h2>

      <form onSubmit={add} className="add-form">
        <input
          className="field-input"
          type="number"
          min="1"
          value={count}
          onChange={(e) => setCount(e.target.value)}
          style={{ flex: 1 }}
        />

        <button className="btn-primary" type="submit">
          Добавить места
        </button>
      </form>

      {loading ? (
        <Spinner />
      ) : (
        <div className="card-grid">
          {slots.map((s) => (
            <div className="vehicle-card" key={s.slot_id}>
              <div className="vehicle-plate">
                Место #{s.slot_id}
              </div>

              <Badge color={s.is_occupied ? "gray" : "green"}>
                {s.is_occupied ? "Занято" : "Свободно"}
              </Badge>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default function App() {
  const [auth, setAuth] = useState(() => {
    const token = localStorage.getItem("token");
    const userId = localStorage.getItem("userId");
    const email = localStorage.getItem("userEmail");

    return token
      ? { token, user_id: userId, email }
      : null;
  });
  const [page, setPage] = useState("vehicles");
  const [toast, setToast] = useState({ msg: "", type: "ok" });
  const [vehicles, setVehicles] = useState([]);
  const [activeSessions, setActiveSessions] = useState([]);
  const [history, setHistory] = useState(() => {
    const saved = localStorage.getItem("parking_history");
    return saved ? JSON.parse(saved) : [];
  });
  const [showModal, setShowModal] = useState(false);
  const [modalVehicle, setModalVehicle] = useState(null);
  useEffect(() => {
    localStorage.setItem("parking_history", JSON.stringify(history));
  }, [history]);

  useEffect(() => {
    function handleUnauthorized() {
      setAuth(null);
      setVehicles([]);
      setActiveSessions([]);
      setHistory([]);
    }

    window.addEventListener("unauthorized", handleUnauthorized);

    return () => {
      window.removeEventListener("unauthorized", handleUnauthorized);
    };
  }, []);



  function showToast(msg, type = "ok") {
    setToast({ msg, type });
  }

  function handleLogout() {
    api.logout().catch(() => {});
    localStorage.clear();
    setAuth(null);
    setActiveSessions([]);
  }

  function handleStartSession(plate) {
    setModalVehicle(plate);
    setShowModal(true);
  }

  function handleSessionStarted(sess) {
    setActiveSessions(prev => [...prev, sess]);
    setPage("sessions");
  }

  function handleSessionEnded(id, data) {
    setActiveSessions(prev => prev.filter(s => s.session_id !== id));

    setHistory(prev => [
      {
        session_id: data.session_id,
        slot_id: data.slot_id,
        vehicle_id: data.vehicle_id,
        start_time_unix: data.start_time_unix,
        end_time_unix: data.end_time_unix,
        amount: data.amount,
      },
      ...prev,
    ]);
  }

  // load vehicles list to pass into modal
  async function refreshVehicles() {
    try {
      const d = await api.vehicles();
      setVehicles(d.vehicles || []);
    } catch {}
  }

  useEffect(() => {
    if (!auth) return;

    refreshVehicles();

    async function loadSessions() {

      try {

        const d =
          await api.activeSessions();

        setActiveSessions(
          d.sessions || []
        );

      } catch (e) {

        console.error(e);
      }
    }

    loadSessions();

  }, [auth]);

  if (!auth) {
    return <AuthPage onAuth={data => { setAuth(data); }} />;
  }

  const navItems = [
    { id: "profile", label: "Кабинет", icon: Icons.user },
    { id: "vehicles", label: "Автомобили", icon: Icons.car },
    { id: "sessions", label: "Парковки", icon: Icons.clock, badge: activeSessions.length || null },
    { id: "history", label: "История", icon: Icons.parking },
    { id: "admin", label: "Админ", icon: Icons.parking },
  ];

  return (
    <div className="app">
      <Toast msg={toast.msg} type={toast.type} onClose={() => setToast({ msg: "" })} />

      <aside className="sidebar">
        <div className="sidebar-logo">
          <div className="logo-icon sm"><Icon d={Icons.parking} size={18} stroke="#fff" /></div>
          <span className="sidebar-brand">SmartPark</span>
        </div>

        <nav className="sidebar-nav">
          {navItems.map(n => (
            <button
              key={n.id}
              className={`nav-item ${page === n.id ? "active" : ""}`}
              onClick={() => setPage(n.id)}
            >
              <Icon d={n.icon} size={18} />
              <span>{n.label}</span>
              {n.badge ? <span className="nav-badge">{n.badge}</span> : null}
            </button>
          ))}
        </nav>

        <div className="sidebar-bottom">
          <SlotsWidget />
          <button className="nav-item logout" onClick={handleLogout}>
            <Icon d={Icons.logout} size={18} />
            <span>Выйти</span>
          </button>
        </div>
      </aside>

      <main className="main-content">
        {page === "profile" && (
          <ProfilePage
            auth={auth}
            vehicles={vehicles}
            activeSessions={activeSessions}
            history={history}
          />
        )}
        {page === "vehicles" && (
          <VehiclesPage
            onStartSession={handleStartSession}
            toast={showToast}
          />
        )}
        {page === "sessions" && (
          <SessionsPage
            activeSessions={activeSessions}
            onEnd={handleSessionEnded}
            toast={showToast}
          />
        )}
        {page === "history" && <HistoryPage history={history} />}
        {page === "admin" && (
          <AdminSlotsPage toast={showToast} />
        )}
      </main>

      {showModal && (
        <StartSessionModal
          vehicleId={modalVehicle}
          vehicles={vehicles}
          onClose={() => { setShowModal(false); setModalVehicle(null); }}
          onStarted={handleSessionStarted}
          toast={showToast}
        />
      )}
    </div>
  );
}
