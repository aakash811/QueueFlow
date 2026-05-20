"use client";

import { useEffect, useState } from "react";

import { api } from "@/lib/api";

interface Job {
  id: string;
  queue_name: string;
  retry_count: number;
}

export default function DeadLetterPage() {
  const [jobs, setJobs] = useState<Job[]>([]);

  useEffect(() => {
    async function fetchJobs() {
      const res = await api.get("/dead-letter/recent");

      setJobs(res.data);
    }

    fetchJobs();
  }, []);

  return (
    <div>
      <h1 className="text-4xl font-bold mb-8">Dead Letter Queue</h1>

      <div className="bg-zinc-900 rounded-2xl border border-zinc-800 overflow-hidden">
        <table className="w-full">
          <thead className="bg-zinc-800">
            <tr>
              <th className="text-left p-4">Job ID</th>

              <th className="text-left p-4">Queue</th>

              <th className="text-left p-4">Retries</th>
            </tr>
          </thead>

          <tbody>
            {jobs.map((job) => (
              <tr key={job.id} className="border-t border-zinc-800">
                <td className="p-4 text-sm">{job.id}</td>

                <td className="p-4">{job.queue_name}</td>

                <td className="p-4 text-red-400">{job.retry_count}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
