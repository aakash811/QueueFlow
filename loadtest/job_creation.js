import http from "k6/http";

import { sleep } from "k6";

export const options = {
  vus: 20,
  duration: "30s",
};

export default function () {
  const loginPayload = JSON.stringify({
    username: "admin",
    password: "password",
  });

  const loginRes = http.post("http://localhost:8080/login", loginPayload, {
    headers: {
      "Content-Type": "application/json",
    },
  });

  const token = loginRes.json("token");

  const payload = JSON.stringify({
    queue_name: "email",
    payload: {
      fail: true,
    },
  });

  const headers = {
    Authorization: `Bearer ${token}`,
    "Idempotency-Key": `${__VU}-${__ITER}`,
    "Content-Type": "application/json",
  };

  http.post("http://localhost:8080/jobs", payload, { headers });

  sleep(1);
}
