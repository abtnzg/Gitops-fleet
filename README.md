# gitops-fleet

This repo declares the **clusters** and **app-of-apps** for the platform.

It is the entry point that ArgoCD watches. Everything in `apps/` is rendered
as ArgoCD `Application` resources; everything in `clusters/` is rendered as
`AppProject` + per-cluster overrides.

## Layout

```
gitops-fleet/
├── apps/                          # App-of-Apps: one Application per workload
│   ├── platform-tooling/          # Argo Workflows, ArgoCD, cert-manager, ...
│   │   ├── base/
│   │   └── overlays/
│   │       ├── dev/
│   │       ├── staging/
│   │       └── prod/
│   ├── observability/             # Mimir, Loki, Tempo, Grafana
│   └── argocd-image-updater/      # Image promotion automation
├── clusters/                      # Cluster-scoped resources
│   ├── mgmt-aws-use1/             # management cluster (us-east-1)
│   │   ├── argocd/                # ArgoCD bootstrap
│   │   ├── policies/              # Gatekeeper constraints
│   │   ├── projects/              # AppProjects (per-team RBAC)
│   │   └── secrets/               # ExternalSecret manifests (no values!)
│   ├── mgmt-aws-euw1/             # (phase 2) management cluster (eu-west-1)
│   ├── mgmt-azure-eastus/         # (phase 3) management cluster (Azure)
│   ├── mgmt-gcp-uscentral1/       # (phase 3) management cluster (GCP)
│   └── spoke-*/                   # workload clusters (phase 2+)
├── projects/                      # Shared AppProject templates
│   └── team-app-project.yaml      # default per-team AppProject
└── README.md
```

## App-of-Apps pattern

The root `Application` declared in the management cluster's bootstrap
(see `platform-fleet/scripts/bootstrap.sh`) points at `apps/`. Each entry
under `apps/<name>/overlays/<env>/` is itself an ArgoCD `Application` that
references a Helm chart, Kustomize directory, or git URL.

Adding a new platform component:
1. Create `apps/<component>/base/` with manifests or a kustomization.
2. Create `apps/<component>/overlays/dev|staging|prod/` with environment overlays.
3. Commit. The root app-of-apps picks it up automatically.

## AppProject model

`AppProject` is the RBAC primitive in ArgoCD. We use one per team:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: team-platform
  namespace: argocd
spec:
  description: Platform team
  sourceRepos:
    - 'https://github.com/your-org/*'
  destinations:
    - namespace: 'team-platform-*'
      server: '*'
  clusterResourceWhitelist:
    - group: ''
      kind: Namespace
```

Teams get their own AppProject. They cannot write to other namespaces
or repos outside their source allowlist.
