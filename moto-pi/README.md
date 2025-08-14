
| Module      | Responsibilities                                                                                  |
| ----------- | ------------------------------------------------------------------------------------------------- |
| HAL         | Provide low-level interfaces for controlling hardware, shielding other parts from hardware quirks |
| Network     | Scan, connect, monitor WiFi; manage Bluetooth pairing; control cellular modem and failovers       |
| Security    | Enforce authentication, encrypt sensitive data, configure VPN and firewall rules                  |
| Services    | Run logic loops for relays, GPS reading, RFID processing, modem status, emergency triggers        |
| Frontend    | Present data & control UI, handle user input, display real-time device status                     |
| Backend API | Secure API endpoints, data validation, interface between UI and firmware, store persistent data   |


Power Resiliency
Use journaling filesystem (e.g., ext4 with journaling) to avoid corruption on power loss.

Implement graceful shutdown triggers via GPIO button or relay detection.

Enable watchdog timers on the Pi to auto-reboot on hangs or crashes.

Use atomic writes (write temp + rename) for critical config and logs.

Save important state in non-volatile storage ASAP (e.g., SQLite with WAL mode).

Network Resiliency
Network Manager Script continuously:

Prioritizes WiFi > Cellular > Bluetooth for connectivity.

Automatically reconnects on disconnect.

Performs periodic health checks (ping a known server).

Implements failover:

If WiFi down, switch to cellular modem automatically.

If both down, enable Bluetooth tethering (if paired device available).

Use Tailscale or WireGuard VPN for secure remote access, reconnect automatically.

Implement local-only fallback mode:

Webserver runs locally even without internet.

UI disables remote features but allows local control.

Use exponential backoff for reconnect attempts to avoid flooding network.
