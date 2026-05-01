# Inference benchmark

Throughput (RPS) and TTFT/ITL (time-to-first-token, inter-token latency) across:

- Backend: Ollama, vLLM, llama.cpp
- Model: qwen3.5-4b, qwen3-coder-30b-a3b, qwen3.6-27b, llama-3.3-70b-q4
- Concurrency: 1, 4, 16, 64, 256
- Hardware: M-series Mac, RTX 4090, L40S, A100, H100

Results inform the SMB/Enterprise tier inference defaults and the cost-vs-quality knee in the routing rules.
