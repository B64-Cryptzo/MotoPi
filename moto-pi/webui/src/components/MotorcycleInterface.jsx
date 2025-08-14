import { useEffect, useState } from "react";

export default function MotorcycleInterface() {
    const [status, setStatus] = useState("loading...");
    const [location, setLocation] = useState("location: 0 0");

    useEffect(() => {
        // Fetch status
        fetch("http://localhost:8080/v1/api/motorcycle/status")
            .then((res) => res.json())
            .then((data) => setStatus(data.status))
            .catch(() => setStatus("error"));

        // Fetch GPS location
        fetch("http://localhost:8080/v1/api/motorcycle/gps")
            .then((res) => res.json())
            .then((data) => setLocation(`location: ${data.lat} ${data.lng}`))
            .catch(() => setLocation("location: error"));
    }, []);

    return (
        <div>
            <h1>Motorcycle Status: {status}</h1>
            <h2>{location}</h2>
        </div>
    );
}