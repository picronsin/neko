#!/usr/bin/env bash
set -euo pipefail

integration_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
output_dir="${NEKO_E2E_OUTPUT_DIR:-${integration_dir}/results/baseline}"

: "${NEKO_E2E_PASSWORD:?set NEKO_E2E_PASSWORD}"

mkdir -p "${output_dir}"

reports=()
for viewers in 1 2 5; do
  run_dir="${output_dir}/viewers-${viewers}"
  mkdir -p "${run_dir}"
  echo "==> running ${viewers}-viewer baseline"
  NEKO_E2E_VIEWERS="${viewers}" \
    NEKO_E2E_OUTPUT_DIR="${run_dir}" \
    "${integration_dir}/benchmark.sh"

  report="${output_dir}/viewers-${viewers}.md"
  "${integration_dir}/report.sh" "${run_dir}" "${report}"
  reports+=("${report}")
done

{
  echo "# Chromium WebRTC 基线矩阵"
  echo
  echo "- 结果目录：\`${output_dir}\`"
  echo "- viewer 数：1、2、5"
  echo
  echo "## 分 viewer 报告"
  echo
  for report in "${reports[@]}"; do
    echo "- [$(basename "${report}")]($(basename "${report}"))"
  done
} >"${output_dir}/baseline-report.md"
echo "baseline report: ${output_dir}/baseline-report.md"
