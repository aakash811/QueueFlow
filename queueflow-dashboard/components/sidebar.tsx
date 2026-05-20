"use client";

import Link from "next/link";
import { LayoutDashboard, ListTodo, Cpu, AlertTriangle } from "lucide-react";

const links = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    name: "Jobs",
    href: "/jobs",
    icon: ListTodo,
  },
  {
    name: "Workers",
    href: "/workers",
    icon: Cpu,
  },
  {
    name: "Dead Letter",
    href: "/dead-letter",
    icon: AlertTriangle,
  },
];

export default function Sidebar() {
  return (
    <div className="w-64 min-h-screen border-r bg-black text-white p-4">
      <h1 className="text-2xl font-bold mb-8">QueueFlow</h1>

      <div className="space-y-3">
        {links.map((link) => {
          const Icon = link.icon;

          return (
            <Link
              key={link.href}
              href={link.href}
              className="flex items-center gap-3 p-3 rounded-lg hover:bg-zinc-800 transition"
            >
              <Icon size={18} />
              <span>{link.name}</span>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
