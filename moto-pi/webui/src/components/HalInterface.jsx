import { useEffect, useState } from "react";

export default function HalInterface() {
    const [status, setStatus] = useState("loading...");

    useEffect(() => {
        fetch("http://localhost:8080/v1/api/hal/status")
            .then((res) => res.json())
            .then((data) => setStatus(data.status))
            .catch(() => setStatus("error"));
    }, []);

    return (
        <div>
            <h1>HAL Status: {status}</h1>
        </div>
    );
}
