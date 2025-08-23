import "./PieMenu.css";

export default function PieMenu({ items = [], onSelect }) {
  const radius = 48; // circle radius
  const center = 50; // center of SVG
  const angleStep = (2 * Math.PI) / items.length;

  // Arc path for each slice
  const createPath = (index) => {
    const startAngle = index * angleStep - Math.PI / 2;
    const endAngle = startAngle + angleStep;

    const x1 = center + radius * Math.cos(startAngle);
    const y1 = center + radius * Math.sin(startAngle);
    const x2 = center + radius * Math.cos(endAngle);
    const y2 = center + radius * Math.sin(endAngle);

    const largeArc = angleStep > Math.PI ? 1 : 0;

    return `M${center},${center} L${x1},${y1} A${radius},${radius} 0 ${largeArc},1 ${x2},${y2} Z`;
  };

  // Label position
  const getLabelPosition = (index) => {
    const midAngle = index * angleStep + angleStep / 2 - Math.PI / 2;
    const r = radius * 0.6;
    const x = center + r * Math.cos(midAngle);
    const y = center + r * Math.sin(midAngle);
    return { x, y };
  };

  return (
    <div className="pie-container">
      <svg viewBox="0 0 100 100" className="pie-svg">
      {items.map((item, i) => {
        const path = createPath(i);
        return (
          <g
            key={i}
            className="pie-slice-group"
            onClick={() => onSelect(item.value)}
          >
            <path d={path} className="pie-slice" />
          </g>
        );
      })}
      </svg>

      {items.map((item, i) => {
        const { x, y } = getLabelPosition(i);
        return (
          <div
            key={i}
            className="pie-label"
            style={{ left: `${x}%`, top: `${y}%` }}
          >
            {/* Decide how to render icon */}
            {typeof item.icon === "string" ? (
              item.icon.startsWith("/") || item.icon.startsWith("http") ? (
                <img src={item.icon} alt={item.value} className="pie-icon" />
              ) : (
                <span className="pie-emoji">{item.icon}</span>
              )
            ) : (
              item.icon || item.label
            )}
          </div>
        );
      })}
    </div>
  );
}
