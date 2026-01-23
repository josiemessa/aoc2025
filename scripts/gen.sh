#!/usr/bin/env bash
set -euxo pipefail

force=0
if [[ "${1:-}" == "--force" ]]; then
  force=1
fi

cwd_base="$(basename "$PWD")"
if [[ "$cwd_base" =~ ^day[0-9]+$ ]]; then
  echo "refusing to run from $PWD; run from repo root instead." >&2
  exit 1
fi

max_day=0
for d in day*; do
  if [[ -d "$d" && "$d" =~ ^day([0-9]+)$ ]]; then
    n="${BASH_REMATCH[1]}"
    if (( n > max_day )); then
      max_day="$n"
    fi
  fi
done

target_dir="day$((max_day + 1))"
mkdir -p "$target_dir"

main_path="$target_dir/main.go"
if [[ -e "$main_path" && "$force" -ne 1 ]]; then
  echo "skip $main_path (exists). use --force to overwrite." >&2
else
  cat > "$main_path" <<'EOF'
package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/josiemessa/aoc2025/pkg/utils"
)

func main() {
	log.SetFlags(0)
	debug := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()
	if !*debug {
		log.SetOutput(io.Discard)
	}

	lines := utils.ReadFileAsLines("input")
	fmt.Println("Lines:", len(lines))
}
EOF
fi

input_path="$target_dir/input"
if [[ -e "$input_path" && "$force" -ne 1 ]]; then
  echo "skip $input_path (exists). use --force to overwrite." >&2
else
  : > "$input_path"
fi

test_input_path="$target_dir/test-input"
if [[ -e "$test_input_path" && "$force" -ne 1 ]]; then
  echo "skip $test_input_path (exists). use --force to overwrite." >&2
else
  : > "$test_input_path"
fi
