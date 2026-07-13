import http from "k6/http";
import { sleep } from "k6";

export const options = {
  scenarios: {
    load_test: {
      executor: "per-vu-iterations",
      vus: 100,
      iterations: 200,
      maxDuration: "2m",
    },
  },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

export function setup() {
  const loginRes = http.post(
    BASE_URL + "/login",
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

export default function (data) {
  const payload = JSON.stringify({
    queue_name: "email",
    payload: {
      message: "benchmark",
    },
  });

  const idemKey = String(__VU) + "-" + String(__ITER) + "-" + String(Date.now());
  const headers = {
    Authorization: "Bearer " + data.token,
    "Idempotency-Key": idemKey,
    "Content-Type": "application/json",
  };

  const res = http.post(BASE_URL + "/jobs", payload, { headers });

  if (res.status !== 201) {
    console.error("Job creation failed:", res.body);
  }
}

export function teardown(data) {
  console.log("Benchmark complete. Check Grafana for metrics.");
}
