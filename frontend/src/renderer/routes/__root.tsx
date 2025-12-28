import {
  createRootRoute,
  Link,
  Outlet,
  useLocation,
} from "@tanstack/react-router";
import {
  BarChart3,
  Bot,
  Circle,
  Info,
  Key,
  LayoutDashboard,
  Settings as SettingsIcon,
  Users,
} from "lucide-react";
import { cn } from "../lib/utils";
import { useProxyStatus } from "../queries";

export const Route = createRootRoute({
  component: Root,
});

function Root() {
  return (
    <div className="flex h-screen overflow-hidden bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/20">
      <Sidebar />
      <main className="flex-1 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  );
}

function Sidebar() {
  const location = useLocation();
  const { data: proxyStatus } = useProxyStatus();

  const navItems = [
    { to: "/", icon: LayoutDashboard, label: "Dashboard" },
    { to: "/quota", icon: BarChart3, label: "Quota" },
    { to: "/providers", icon: Users, label: "Providers" },
    { to: "/agents", icon: Bot, label: "Agents" },
    { to: "/api-keys", icon: Key, label: "API Keys" },
    { to: "/settings", icon: SettingsIcon, label: "Settings" },
    { to: "/about", icon: Info, label: "About" },
  ];

  return (
    <aside className="w-64 bg-white/80 backdrop-blur-sm border-r-2 border-black flex flex-col shadow-lg">
      <div className="p-4 border-b-2 border-black bg-gradient-to-r from-neobrutal-blue/5 to-transparent">
        <div className="flex items-center gap-3">
          <img
            src="/logo.png"
            alt="Quotio Logo"
            className="w-10 h-10 rounded-sm border-2 border-black shadow-neobrutal-sm object-cover"
          />
          <div>
            <h1 className="text-xl font-bold tracking-tight bg-gradient-to-r from-neobrutal-blue to-neobrutal-purple bg-clip-text text-transparent">
              Quotio
            </h1>
            <p className="text-xs text-gray-500">AI Proxy Manager</p>
          </div>
        </div>
      </div>
      <nav className="flex-1 py-3 px-2">
        {navItems.map((item) => {
          const isActive = location.pathname === item.to;
          return (
            <Link
              key={item.to}
              to={item.to}
              className={cn(
                "flex items-center gap-3 px-4 py-2.5 mb-1 font-medium rounded-sm",
                "transition-all duration-200 ease-out group",
                isActive
                  ? "bg-black text-white shadow-neobrutal-sm"
                  : "text-gray-600 hover:bg-gray-100 hover:text-black hover:translate-x-1"
              )}
            >
              <item.icon
                className={cn(
                  "w-5 h-5 transition-transform duration-200",
                  !isActive && "group-hover:scale-110"
                )}
              />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>
      <div className="p-4 border-t-2 border-black bg-gray-50/50">
        <div className="flex items-center gap-2.5 text-sm">
          <div className="relative">
            <Circle
              className={cn(
                "w-3 h-3",
                proxyStatus?.running
                  ? "fill-emerald-500 text-emerald-500"
                  : "fill-red-500 text-red-500"
              )}
            />
            {proxyStatus?.running && (
              <span className="absolute inset-0 w-3 h-3 bg-emerald-400 rounded-full animate-ping opacity-75" />
            )}
          </div>
          <span className="font-semibold">
            {proxyStatus?.running ? "Running" : "Stopped"}
          </span>
          <span className="text-gray-500 font-mono text-xs">
            :{proxyStatus?.port || 8317}
          </span>
        </div>
      </div>
    </aside>
  );
}
