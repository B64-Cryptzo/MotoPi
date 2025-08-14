import { useEffect, useState } from "react";

export default function NetworkScanner() {
    const [status, setStatus] = useState("loading...");

    useEffect(() => {
        fetch("http://localhost:8080/v1/api/network/status")
            .then((res) => res.json())
            .then((data) => setStatus(data.status))
            .catch(() => setStatus("error"));
    }, []);

    return (
        <div>
            <h1>Network Status: {status}</h1>
        </div>
    );
}
