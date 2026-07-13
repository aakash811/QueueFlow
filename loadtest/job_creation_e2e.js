import http from "k6/http";
import { sleep } from "k6";

export const options = {
  scenarios: {
    submission: {
      executor: "per-vu-iterations",
      vus: 50,
      iterations: 100,
      maxDuration: "1m",
      exec: "submitJob",
    },
    verification: {
      executor: "constant-vus",
      vus: 1,
      duration: "2m",
      exec: "verifyCompletion",
      startTime: "30s",
    },
  },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

export function setup() {
  const loginRes = http.post(
    `${BASE_URL}/login`,
    JSON.stringify({ username: "admin", password: "admin123" }),
    {
      headers: { "Content-Type": "application/json" },
    },
  );

  if (loginRes.status !== 200) {
    console.error("Login failed:", loginRes.body);
    return { token: "" };
  }

  const token = loginRes.json("token");
  return { token };
}

export function submitJob(data) {
  const payload = JSON.stringify({
    queue_name: "email",
    payload: {
      message: "benchmark",
    },
  });

  const headers = {
    Authorization: `Bearer ${data.token}`,
    Idempotency-Key: `${__VU}-${__ITER}-${Date.now()}`,
    "Content-Type": "application/json",
  };

  const res = http.post(`${BASE_URL}/jobs`, payload, { headers });

  if (res.status !== 201) {
    console.error("Job creation failed:", res.body);
  }

  sleep(0.01);
}

export function verifyCompletion(data) {
  const res = http.get(`${BASE_URL}/jobs?queue_name=email&limit=5`, {
    headers: { Authorization: `Bearer ${data.token}` },
  });

  if (res.status === 200) {
    const jobs = res.json();
    if (jobs && jobs.length > 0) {
      console.log(`Latest job status: ${jobs[0].status}`);
    }
  }

  sleep(5);
}

export function teardown(data) {
  console.log("Benchmark complete.");
}
