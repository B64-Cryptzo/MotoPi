import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";
import PieMenu from "./components/PieMenu";
import StatusPage from "./pages/status/StatusPage";
import "./App.css";

function MenuWrapper() {
  const navigate = useNavigate();

  const items = [
    { value: "status", label: "Status", icon: "📊" },
    { value: "network", label: "Network", icon: "🌐" },
    { value: "motorcycle", label: "Moto", icon: "🏍️" },
    { value: "settings", label: "Settings", icon: "⚙️" },
  ];

  const handleSelect = (value) => {
    navigate(`/${value}`);
  };

  return (
    <div className="app-container">
      <PieMenu items={items} onSelect={handleSelect} />
    </div>
  );
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MenuWrapper />} />
        <Route path="/status" element={<StatusPage />} />
        {/* Add more pages as needed */}
      </Routes>
    </Router>
  );
}

export default App;
