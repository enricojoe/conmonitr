import { NavLink } from "react-router-dom";

const links = [
  { to: "/", label: "Containers", end: true },
  { to: "/images", label: "Images" },
  { to: "/volumes", label: "Volumes" },
  { to: "/networks", label: "Networks" },
];

export default function NavMenu() {
  return (
    <nav className="mb-6 flex gap-1 border-b border-zinc-800 pb-3">
      {links.map(({ to, label, end }) => (
        <NavLink
          key={to}
          to={to}
          end={end}
          className={({ isActive }) =>
            `rounded px-3 py-1.5 text-sm font-medium transition-colors ${
              isActive
                ? "bg-zinc-700 text-zinc-100"
                : "text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200"
            }`
          }
        >
          {label}
        </NavLink>
      ))}
    </nav>
  );
}
