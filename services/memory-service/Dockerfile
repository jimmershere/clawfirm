FROM python:3.12-slim AS build
WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends build-essential && rm -rf /var/lib/apt/lists/*
COPY pyproject.toml ./
RUN pip install --no-cache-dir .
COPY memory_service ./memory_service
COPY proto ./proto

FROM python:3.12-slim
WORKDIR /app
COPY --from=build /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=build /usr/local/bin /usr/local/bin
COPY --from=build /src/memory_service ./memory_service
COPY --from=build /src/proto ./proto
RUN useradd -r -u 10001 memsvc && chown -R memsvc:memsvc /app
USER memsvc
EXPOSE 50053 8080
ENTRYPOINT ["python", "-m", "memory_service.api"]
