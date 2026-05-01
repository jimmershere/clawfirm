# Enterprise teardown

Enterprise teardown is intentionally manual because it is destructive and often tied to compliance, data-retention, and incident-response workflows.

## Recommended order

1. Export audit logs, approvals, and policy bundles from ClawSecure.
2. Snapshot OpenBao, Postgres, Langfuse, and object storage.
3. Suspend Flux reconciliation for the cluster.
4. Drain workloads and confirm no active approval queues or long-running jobs remain.
5. Remove external DNS, ingress, and zero-trust access paths.
6. Tear down worker nodes, then control-plane nodes, through your normal k0sctl and infrastructure workflows.
7. Revoke signing keys, OIDC clients, service accounts, and overlay-network credentials.

## Notes

- Treat teardown as a change-controlled operation.
- If the cluster handled regulated data, preserve chain-of-custody records before deletion.
- Prefer destroying infrastructure only after verified backups and export checks complete.
