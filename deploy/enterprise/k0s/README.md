# Enterprise tier — k0s control plane

Use [`k0sctl`](https://docs.k0sproject.io/main/k0sctl-install/) to bring up an HA k0s cluster from `k0sctl.yaml.example`.

```bash
cp k0sctl.yaml.example k0sctl.yaml
# edit hostnames, SSH keys, labels
k0sctl apply --config k0sctl.yaml
k0sctl kubeconfig --config k0sctl.yaml > ~/.kube/clawfirm-enterprise
export KUBECONFIG=~/.kube/clawfirm-enterprise
kubectl get nodes
```

After the cluster is up, bootstrap Flux against the GitOps repo (see `../flux/`) and Flux will reconcile the rest of the stack from `../helm/`.
