import React from "react";

export default function StatusCard({ title, onClick }) {
  return (
    <div className="status-card" onClick={onClick}>
      <h2 className="status-card-title">{title}</h2>
    </div>
  );
}
