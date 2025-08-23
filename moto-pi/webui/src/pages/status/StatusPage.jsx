import { useState } from "react";
import StatusCard from "../../components/StatusCard";
import "./StatusPage.css";

export default function StatusPage() {
  const [modalData, setModalData] = useState(null);

  const openModal = (title, fetchUrl) => {
    const modalState = { title, statusMap: null, loading: true, error: false };
    setModalData(modalState);

    fetch(fetchUrl)
      .then((res) => res.json())
      .then((data) => setModalData({ ...modalState, statusMap: data, loading: false }))
      .catch(() => setModalData({ ...modalState, error: true, loading: false }));
  };

  const closeModal = () => setModalData(null);

  return (
    <div className="status-page">
      <button className="back-btn" onClick={() => window.history.back()}>
        ⬅ Back
      </button>

      <h1 className="page-title">System Status</h1>

      <div className="status-grid">
        <StatusCard title="I/O" onClick={() => openModal("I/O", "http://10.10.10.1:8080/v1/api/hal/status")} />
        <StatusCard title="Network" onClick={() => openModal("Network", "http://10.10.10.1:8080/v1/api/network/status")} />
        <StatusCard title="Motorcycle" onClick={() => openModal("Motorcycle", "http://10.10.10.1:8080/v1/api/motorcycle/status")} />
      </div>

      {modalData && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>{modalData.title} Components</h2>
            {modalData.loading && <p>Loading...</p>}
            {modalData.error && <p className="text-red-500">Error fetching status</p>}
            {modalData.statusMap && (
              <ul className="status-detail-list">
                {Object.entries(modalData.statusMap).map(([device, status]) => (
                  <li key={device} className="status-detail-item">
                    <span className="device-name">{device}</span>
                    <span className={status === "online" ? "status-online" : "status-offline"}>
                      {status.toUpperCase()}
                    </span>
                  </li>
                ))}
              </ul>
            )}
            <button className="modal-close-btn" onClick={closeModal}>Close</button>
          </div>
        </div>
      )}
    </div>
  );
}
