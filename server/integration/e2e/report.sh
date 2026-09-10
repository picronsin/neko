#!/usr/bin/env bash
set -euo pipefail

results_dir="${1:-${NEKO_E2E_OUTPUT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/results}}"
output="${2:-${results_dir}/baseline-report.md}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required to generate the baseline report. Install it with: sudo apt-get update && sudo apt-get install -y jq" >&2
  exit 2
fi

shopt -s nullglob
results=("${results_dir}"/viewer-*.json)
metrics=("${results_dir}"/viewer-*.prom)
if [[ "${#results[@]}" -eq 0 ]]; then
  echo "no viewer JSON results found in ${results_dir}" >&2
  exit 1
fi

mkdir -p "$(dirname "${output}")"
{
  echo "# Chromium WebRTC 基线报告"
  echo
  echo "- 结果目录：\`${results_dir}\`"
  echo "- 生成时间：\`$(date -u +%Y-%m-%dT%H:%M:%SZ)\`"
  echo
  echo "## 观看者结果"
  echo
  echo "| 样本 | Profile | 视口 | 视频 | 连接耗时 | 首帧耗时 | Envelope | 状态 |"
  echo "| --- | --- | --- | --- | ---: | ---: | ---: | --- |"
  for result in "${results[@]}"; do
    file="$(basename "${result}")"
    jq -r --arg file "${file}" '
      if .error then
        "| \($file) | \(.profile // "-") | - | - | - | - | \(.malformedFrames // "-") | FAILED: \(.error) |"
      else
        "| \($file) | \(.profile // "-") | \(.viewport.width)x\(.viewport.height) | \(.video.width)x\(.video.height) | \(.connectionMs) ms | \(.firstFrameMs) ms | \(.malformedFrames) malformed | PASS |"
      end
    ' "${result}"
  done

  summary="$(jq -s -r '
    map(select(.error == null)) as $ok |
    map(select(.error != null)) as $failed |
    [$ok[]?.connectionMs] as $connections |
    [$ok[]?.firstFrameMs] as $firstFrames |
    def average($values): if ($values | length) == 0 then "-" else ((($values | add) / ($values | length)) | round | tostring) end;
    def percentile($values; $p):
      ($values | sort) as $sorted |
      ($sorted | length) as $count |
      if $count == 0 then "-" else ($sorted[((($count - 1) * $p) | floor)] | tostring) end;
    "- 通过样本：\($ok | length)"
    + "\n- 失败样本：\($failed | length)"
    + "\n- 连接耗时平均值：\(average($connections)) ms"
    + "\n- 连接耗时 P95：\(percentile($connections; 0.95)) ms"
    + "\n- 首帧耗时平均值：\(average($firstFrames)) ms"
    + "\n- 首帧耗时 P95：\(percentile($firstFrames; 0.95)) ms"
  ' "${results[@]}")"

  echo
  echo "## 汇总"
  echo
  printf '%s\n' "${summary}"

  echo
  echo "## Prometheus 快照"
  echo
  if [[ "${#metrics[@]}" -eq 0 ]]; then
    echo "未找到 Prometheus 快照。"
  else
    for metric in "${metrics[@]}"; do
      echo "- \`$(basename "${metric}")\`"
    done
  fi
} >"${output}"
