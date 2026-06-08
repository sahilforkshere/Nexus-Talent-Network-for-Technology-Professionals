#!/bin/bash
# Post 30 diverse test jobs to Nexus via GraphQL
# Usage: TOKEN=<jwt> ./scripts/seed_jobs.sh

set -e

if [ -z "$TOKEN" ]; then
  echo "Set TOKEN env var first: export TOKEN=<your_jwt>"
  exit 1
fi

URL="http://localhost:4000/"

post_job() {
  local title="$1"
  local company="$2"
  local location="$3"
  local type="$4"
  local level="$5"
  local desc="$6"

  local payload
  payload=$(jq -n \
    --arg t "$title" --arg c "$company" --arg l "$location" \
    --arg jt "$type" --arg el "$level" --arg d "$desc" \
    '{
      query: "mutation PostJob($input: PostJobInput!) { postJob(input: $input) { job_id title company } }",
      variables: { input: { title: $t, company: $c, location: $l, job_type: $jt, experience_level: $el, description: $d } }
    }')

  curl -s -X POST "$URL" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$payload" \
    | jq -r 'if .data.postJob then "\(.data.postJob.job_id) | \(.data.postJob.title) @ \(.data.postJob.company)" else "ERROR: \(.errors[0].message)" end'
}

echo "Posting 30 jobs..."

post_job "Senior Backend Engineer" "Google" "Mountain View, CA" "FULL_TIME" "SENIOR" "Build large-scale distributed systems in Go and C++. Work on Bigtable and Spanner infrastructure."
post_job "Staff Software Engineer" "Meta" "Menlo Park, CA" "FULL_TIME" "LEAD" "Design microservices architecture for social graph. Experience with Thrift and distributed databases required."
post_job "Platform Engineer" "Stripe" "San Francisco, CA" "FULL_TIME" "MID" "Build payment infrastructure handling millions of transactions. Go, Kafka, Postgres."
post_job "Infrastructure Engineer" "Cloudflare" "Remote" "FULL_TIME" "MID" "Work on global CDN infrastructure. Rust and Go preferred. Distributed systems experience essential."
post_job "Site Reliability Engineer" "Netflix" "Los Gatos, CA" "FULL_TIME" "SENIOR" "Ensure reliability of streaming platform serving 200M users. Kubernetes, Prometheus, Go."
post_job "Machine Learning Engineer" "OpenAI" "San Francisco, CA" "FULL_TIME" "SENIOR" "Train large language models. PyTorch, CUDA, distributed training on GPU clusters."
post_job "AI Research Engineer" "DeepMind" "London, UK" "FULL_TIME" "SENIOR" "Implement novel deep learning architectures. Research background in transformers preferred."
post_job "Data Scientist" "Airbnb" "San Francisco, CA" "FULL_TIME" "MID" "Analyze booking patterns and build recommendation models. Python, Spark, SQL."
post_job "MLOps Engineer" "Databricks" "Remote" "FULL_TIME" "MID" "Build ML pipelines and model serving infrastructure. MLflow, Kubernetes, Python."
post_job "Computer Vision Engineer" "Tesla" "Palo Alto, CA" "FULL_TIME" "SENIOR" "Develop perception systems for autonomous vehicles. C++, CUDA, real-time inference."
post_job "Senior Frontend Engineer" "Figma" "San Francisco, CA" "FULL_TIME" "SENIOR" "Build collaborative design tools in TypeScript and WebAssembly. Performance-critical UI."
post_job "React Developer" "Shopify" "Remote" "FULL_TIME" "MID" "Build merchant-facing dashboards. React, GraphQL, TypeScript. E-commerce domain knowledge."
post_job "Full Stack Engineer" "Linear" "Remote" "FULL_TIME" "MID" "Build fast project management tooling. React, Node.js, PostgreSQL. Focus on performance."
post_job "iOS Engineer" "Spotify" "Stockholm, Sweden" "FULL_TIME" "SENIOR" "Build audio streaming features for iOS app. Swift, AVFoundation, background audio."
post_job "Android Engineer" "Grab" "Singapore" "FULL_TIME" "MID" "Develop ride-hailing features for Android. Kotlin, Jetpack Compose, location services."
post_job "DevOps Engineer" "HashiCorp" "Remote" "FULL_TIME" "MID" "Build Terraform providers and improve CI/CD pipelines. Go, AWS, Kubernetes."
post_job "Cloud Architect" "AWS" "Seattle, WA" "FULL_TIME" "LEAD" "Design cloud solutions for enterprise customers. Deep AWS expertise, solution design."
post_job "Kubernetes Engineer" "Red Hat" "Remote" "FULL_TIME" "MID" "Contribute to OpenShift and upstream Kubernetes. Go, container runtimes, networking."
post_job "Security Engineer" "GitHub" "Remote" "FULL_TIME" "SENIOR" "Protect developer infrastructure at scale. AppSec, threat modeling, Go, Python."
post_job "Database Engineer" "PlanetScale" "Remote" "FULL_TIME" "SENIOR" "Work on distributed MySQL at global scale. Vitess, Go, sharding, replication."
post_job "Golang Developer" "Temporal" "Remote" "FULL_TIME" "MID" "Build workflow orchestration platform in Go. Distributed systems, gRPC, persistence layers."
post_job "Backend Software Engineer" "Vercel" "Remote" "FULL_TIME" "MID" "Build edge compute and deployment infrastructure. Node.js, Go, serverless, CDN."
post_job "Systems Programmer" "Oxide Computer" "San Jose, CA" "FULL_TIME" "SENIOR" "Write firmware and OS-level code in Rust. Low-level systems, hardware integration."
post_job "API Engineer" "Twilio" "Remote" "FULL_TIME" "MID" "Design and build communication APIs. REST, webhooks, Golang, high availability."
post_job "Search Engineer" "Elastic" "Remote" "FULL_TIME" "SENIOR" "Improve Elasticsearch relevance and performance. Java, Lucene, distributed search."
post_job "Junior Backend Developer" "Razorpay" "Bangalore, India" "FULL_TIME" "JUNIOR" "Build payment APIs and internal tools. Java or Go. Fresh graduates welcome."
post_job "Software Engineer Intern" "Microsoft" "Hyderabad, India" "INTERNSHIP" "INTERN" "6-month internship on Azure cloud services. C# or Go. B.Tech final year students."
post_job "Graduate Software Engineer" "Thoughtworks" "Remote" "FULL_TIME" "JUNIOR" "Agile consulting and delivery. Any backend language. Problem solving and learning mindset."
post_job "Associate Engineer" "Zepto" "Mumbai, India" "FULL_TIME" "JUNIOR" "Build quick-commerce backend services. Go, Redis, Postgres. Fast-paced startup."
post_job "Backend Intern" "Cred" "Bangalore, India" "INTERNSHIP" "INTERN" "Work on fintech APIs and data pipelines. Go or Python. IIIT/NIT students preferred."

echo ""
echo "Done! All 30 jobs posted."
