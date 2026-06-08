import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 50,
  duration: "30s",
};

const TOKEN = __ENV.TOKEN;
const URL = "http://localhost:4000/";

const queries = [
  "backend engineer",
  "golang developer",
  "machine learning",
  "mobile iOS Android",
  "devops kubernetes",
];

export default function () {
  const query = queries[Math.floor(Math.random() * queries.length)];

  const res = http.post(
    URL,
    JSON.stringify({
      query: `{ search(query: "${query}") { ... on JobResult { job_id title company } } }`,
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
