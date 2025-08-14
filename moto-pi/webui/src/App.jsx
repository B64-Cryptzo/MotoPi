import './App.css'
import NetworkScanner from "./components/NetworkScanner";
import HalInterface from "./components/HalInterface";
import MotorcycleInterface from "./components/MotorcycleInterface";

function App() {
  return (
        <div>
            <HalInterface />
            <NetworkScanner />
            <MotorcycleInterface />
        </div>
    );
}

export default App
