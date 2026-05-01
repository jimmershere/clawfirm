#!/usr/bin/env bats
# E2E sanity check: bring up the Solo stack and verify all endpoints respond.

setup() {
  export CLAWFIRM_DIR="${BATS_TEST_DIRNAME}/../.."
}

@test "clawfirm CLI builds" {
  run make -C "$CLAWFIRM_DIR" build
  [ "$status" -eq 0 ]
  [ -x "$CLAWFIRM_DIR/bin/clawfirm" ]
}

@test "Solo stack comes up" {
  run make -C "$CLAWFIRM_DIR" solo-up
  [ "$status" -eq 0 ]
  sleep 30
}

@test "ClawSecure dashboard responds" {
  run curl -fsS http://127.0.0.1:3188/api/health
  [ "$status" -eq 0 ]
}

@test "Ollama responds" {
  run curl -fsS http://127.0.0.1:11434/api/version
  [ "$status" -eq 0 ]
}

@test "Solo stack tears down cleanly" {
  run make -C "$CLAWFIRM_DIR" solo-down
  [ "$status" -eq 0 ]
}
