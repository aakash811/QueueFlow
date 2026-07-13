#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-2m}"
RESULTS_DIR="${RESULTS_DIR:-loadtest/results}"
PYTHON_SCRIPT="loadtest/extract_metrics.py"

mkdir -p "$RESULTS_DIR"

echo "QueueFlow Quick Load Test"
echo "========================="
echo "Base URL: $BASE_URL"
echo "Duration: $DURATION"
echo ""

# Check dependencies
if ! command -v k6 &> /dev/null; then
    echo "ERROR: k6 not found. Install from https://k6.io/docs/getting-started/installation/"
    exit 1
fi

if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
    echo "ERROR: Python not found. Install Python 3 to run metrics extraction."
    exit 1
fi

PYTHON_CMD="python3"
if ! command -v python3 &> /dev/null; then
    PYTHON_CMD="python"
fi

# Generate k6 script
cat > loadtest/run.js <<'K6EOF'
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

export default function (data) {
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

export function teardown(data) {
  console.log("Benchmark complete.");
}
K6EOF

# Run k6
echo "Running k6 load test..."
k6 run --out json="$RESULTS_DIR/k6-results.json" loadtest/run.js

echo ""
echo "Waiting 10s for metrics to settle..."
sleep 10

# Collect metrics
echo "Collecting Prometheus metrics..."
$PYTHON_CMD "$PYTHON_SCRIPT" "http://localhost:9090" > "$RESULTS_DIR/metrics.txt"

echo ""
echo "Results:"
cat "$RESULTS_DIR/metrics.txt"
echo ""
echo "Results saved to $RESULTS_DIR"
