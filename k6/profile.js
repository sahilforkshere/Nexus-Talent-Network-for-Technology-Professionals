import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 100,
  duration: "30s",
};

const TOKEN = __ENV.TOKEN;
const URL = "http://localhost:4000/";

export default function () {
  const res = http.post(
    URL,
    JSON.stringify({
      query: `{ listJobs { job_id title company location } }`,
    }),
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${TOKEN}`,
      },
    }
  );

  check(res, {
    "status 200": (r) => r.status === 200,
    "no errors": (r) => !JSON.parse(r.body).errors,
  });

  sleep(0.5);
}
