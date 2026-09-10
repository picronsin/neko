#!/usr/bin/env bash
set -euo pipefail

integration_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
output_dir="${NEKO_E2E_OUTPUT_DIR:-${integration_dir}/results/baseline}"

: "${NEKO_E2E_PASSWORD:?set NEKO_E2E_PASSWORD}"

mkdir -p "${output_dir}"

for viewers in 1 2 5; do
  echo "==> running ${viewers}-viewer baseline"
  NEKO_E2E_VIEWERS="${viewers}" \
    NEKO_E2E_OUTPUT_DIR="${output_dir}" \
    "${integration_dir}/benchmark.sh"
done

"${integration_dir}/report.sh" "${output_dir}"
echo "baseline report: ${output_dir}/baseline-report.md"
