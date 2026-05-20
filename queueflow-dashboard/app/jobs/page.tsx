"use client";

import { useEffect, useState } from "react";

import { api } from "@/lib/api";

interface Job {
  id: string;
  queue_name: string;
  status: string;
  retry_count: number;
}

export default function JobsPage() {
  const [jobs, setJobs] = useState<Job[]>([]);

  useEffect(() => {
    async function fetchJobs() {
      const res = await api.get("/jobs/recent");

      setJobs(res.data);
    }

    fetchJobs();
  }, []);

  return (
    <div>
      <h1 className="text-4xl font-bold mb-8">Jobs</h1>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        <table className="w-full">
          <thead className="bg-zinc-800">
            <tr>
              <th className="text-left p-4">Job ID</th>

              <th className="text-left p-4">Queue</th>

              <th className="text-left p-4">Status</th>

              <th className="text-left p-4">Retries</th>
            </tr>
          </thead>

          <tbody>
            {jobs.map((job) => (
              <tr key={job.id} className="border-t border-zinc-800">
                <td className="p-4 text-sm">{job.id}</td>

                <td className="p-4">{job.queue_name}</td>

                <td className="p-4">
                  <span
                    className={
                      job.status === "completed"
                        ? "text-green-400"
                        : job.status === "failed"
                          ? "text-red-400"
                          : "text-yellow-400"
                    }
                  >
                    {job.status}
                  </span>
                </td>

                <td className="p-4">{job.retry_count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
