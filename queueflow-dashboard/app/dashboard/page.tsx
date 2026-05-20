"use client";

import { useEffect, useState } from "react";

import MetricCard from "@/components/metric-card";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

import { api } from "@/lib/api";
import GrafanaPanel from "@/components/grafana-panel";

export default function DashboardPage() {
  const [metrics, setMetrics] = useState<any>(null);

  useEffect(() => {
    async function fetchMetrics() {
      const res = await api.get("/dashboard/summary");

      setMetrics(res.data);
    }

    fetchMetrics();
  }, []);

  if (!metrics) {
    return <div>Loading...</div>;
  }

  const chartData = [
    {
      name: "Throughput",
      value: metrics.throughput,
    },
    {
      name: "Retries",
      value: metrics.retries,
    },
    {
      name: "Failures",
      value: metrics.failures,
    },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold">QueueFlow Dashboard</h1>

        <p className="text-zinc-400 mt-2">
          Distributed Queue Monitoring Platform
        </p>
      </div>

      <div className="grid grid-cols-3 gap-6">
        <MetricCard title="Queue Depth" value={metrics.queue_depth} />

        <MetricCard title="Throughput" value={metrics.throughput} />

        <MetricCard title="Retries" value={metrics.retries} />

        <MetricCard title="Failures" value={metrics.failures} />

        <MetricCard title="Workers" value={metrics.workers} />

        <MetricCard title="Latency" value={`${metrics.latency}ms`} />
      </div>

      <div className="grid grid-cols-1 gap-6">
        <GrafanaPanel
          title="Queue Throughput"
          url="http://localhost:3000/d-solo/adjz7vr/queueflow-monitoring?orgId=1&theme=dark&panelId=1"
        />

        <GrafanaPanel
          title="Queue Depth"
          url="http://localhost:3000/d-solo/adjz7vr/queueflow-monitoring?orgId=1&theme=dark&panelId=2"
        />

        <GrafanaPanel
          title="Retry Metrics"
          url="http://localhost:3000/d-solo/adjz7vr/queueflow-monitoring?orgId=1&theme=dark&panelId=3"
        />
      </div>
    </div>
  );
}
