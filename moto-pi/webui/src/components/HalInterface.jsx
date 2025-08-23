import { useEffect, useState } from "react";

export default function HalInterface() {
  const [statusMap, setStatusMap] = useState(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/v1/api/hal/status")
      .then((res) => res.json())
      .then((data) => setStatusMap(data))
      .catch(() => setError(true));
  }, []);

  if (error) return <h1>Error fetching status</h1>;
  if (!statusMap) return <h1>Loading...</h1>;

  return (
    <div className="p-4 space-y-2">
      <h1 className="text-xl font-bold">HAL Status</h1>
      <ul className="space-y-1">
        {Object.entries(statusMap).map(([key, value]) => (
          <li key={key} className="flex items-center gap-2">
            <span className="font-semibold capitalize">{key}:</span>
            <span
              className={
                value === "online" ? "text-green-600" : "text-red-600"
              }
            >
              {value}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}
