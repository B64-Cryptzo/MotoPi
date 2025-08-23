import { useState } from "react";
import "./MotorcyclePage.css";

export default function MotorcyclePage() {
  const [modalData, setModalData] = useState(null);

  const triggerAction = (action, url) => {
    const modalState = { action, loading: true, success: null, error: null };
    setModalData(modalState);

    fetch(url, { method: "POST" })
      .then((res) => {
        if (!res.ok) throw new Error("Failed");
        return res.json();
      })
      .then((data) =>
        setModalData({
          ...modalState,
          loading: false,
          success: data.message || "Success",
        })
      )
      .catch(() =>
        setModalData({
          ...modalState,
          loading: false,
          error: "Something went wrong",
        })
      );
  };

  const closeModal = () => setModalData(null);

  return (
    <div className="motorcycle-page">
      <button className="back-btn" onClick={() => window.history.back()}>
        ⬅ Back
      </button>

      <h1 className="page-title">Motorcycle Control</h1>

      <div className="action-grid">
        <div
          className="action-card action-reboot"
          onClick={() =>
            triggerAction(
              "Reboot",
              "http://10.10.10.1:8080/v1/api/motorcycle/reboot"
            )
          }
        >
          Reboot
        </div>
        <div
          className="action-card action-unlock"
          onClick={() =>
            triggerAction(
              "Unlock",
              "http://10.10.10.1:8080/v1/api/motorcycle/unlock"
            )
          }
        >
          Unlock
        </div>
        <div
          className="action-card action-start"
          onClick={() =>
            triggerAction(
              "Start",
              "http://10.10.10.1:8080/v1/api/motorcycle/start"
            )
          }
        >
          Start
        </div>
      </div>

      {modalData && (
        <div className="modal-overlay" onClick={closeModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>{modalData.action} Motorcycle</h2>
            {modalData.loading && <p>Processing...</p>}
            {modalData.success && (
              <p className="status-online">{modalData.success}</p>
            )}
            {modalData.error && (
              <p className="status-offline">{modalData.error}</p>
            )}
            <button className="modal-close-btn" onClick={closeModal}>
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
