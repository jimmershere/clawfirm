FROM python:3.12-slim AS build
WORKDIR /src
COPY pyproject.toml ./
RUN pip install --no-cache-dir .
COPY policy_judge ./policy_judge

FROM python:3.12-slim
WORKDIR /app
COPY --from=build /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=build /usr/local/bin /usr/local/bin
COPY --from=build /src/policy_judge ./policy_judge
RUN useradd -r -u 10001 judge && chown -R judge:judge /app
USER judge
EXPOSE 8080
ENTRYPOINT ["python", "-m", "policy_judge.server"]
