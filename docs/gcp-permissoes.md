# Permissões no GCP (Secret Manager) — TIMGCPSMOPERATOR

## O cluster já tem acesso ao Secret Manager?

**Não por defeito.** Instalar o operador com `kubectl apply` só cria recursos no Kubernetes. **Não** configura IAM na GCP nem liga o pod a uma identidade com permissão para `AccessSecretVersion`.

Sem isso, o processo usa [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials) e **falha** (ou usa uma identidade errada) se não houver credenciais válidas no ambiente do pod.

Em **GKE**, o padrão recomendado é [**Workload Identity**](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity): o `ServiceAccount` do Kubernetes usado pelo Deployment mapeia para uma **Google Service Account (GSA)**; na GCP concedes **só leitura** ao Secret Manager nessa GSA.

---

## O que o operador precisa na IAM

| Necessidade | Papel / permissão típica |
|-------------|---------------------------|
| Ler versões de segredos no GSM | `roles/secretmanager.secretAccessor` |

Aplica o papel **o mais restrito possível**:

- **Por segredo** (recomendado): `gcloud secrets add-iam-policy-binding` só nos `secretId` usados pelos `TimGcpSmSecret`.
- **Por projeto** (mais simples, menos restrito): o mesmo papel ao nível do projeto na GSA (só em ambientes onde aceitares essa superfície).

O código **não** usa APIs de criar/atualizar/apagar segredos no GSM.

---

## Configuração em GKE com Workload Identity

Substitui `PROJECT_ID`, `GCP_SA_NAME` e ajusta o nome da GSA se quiseres.

### 1. Ativar Workload Identity no cluster (se ainda não estiver)

Clusters GKE novos costumam já vir com WI; em dúvida, vê a [documentação oficial](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity#enable).

### 2. Criar a Google Service Account (GSA)

```bash
gcloud iam service-accounts create GCP_SA_NAME \
  --project=PROJECT_ID \
  --display-name="TIMGCPSMOPERATOR controller"
```

### 3. Conceder leitura aos segredos necessários

**Opção A — um segredo específico** (exemplo):

```bash
gcloud secrets add-iam-policy-binding NOME_DO_SEGREDO \
  --project=PROJECT_ID \
  --member="serviceAccount:GCP_SA_NAME@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

Repete por cada segredo que os CRs vão referenciar, ou usa **IAM condicional** / conjunto de segredos conforme a tua política.

**Opção B — ao nível do projeto** (menos restritivo):

```bash
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="serviceAccount:GCP_SA_NAME@PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### 4. Ligar a GSA ao ServiceAccount do Kubernetes

O manifest usa o KSA `timgcpsm-operator` no namespace `timgcpsm-operator-system` ([`config/rbac/service_account.yaml`](../config/rbac/service_account.yaml)).

**Anotação no KSA** (Workload Identity v2 no GKE):

```bash
kubectl annotate serviceaccount timgcpsm-operator \
  -n timgcpsm-operator-system \
  iam.gke.io/gcp-service-account=GCP_SA_NAME@PROJECT_ID.iam.gserviceaccount.com \
  --overwrite
```

**Binding IAM** (permite que o KSA use a GSA):

```bash
gcloud iam service-accounts add-iam-policy-binding \
  GCP_SA_NAME@PROJECT_ID.iam.gserviceaccount.com \
  --project=PROJECT_ID \
  --role="roles/iam.workloadIdentityUser" \
  --member="serviceAccount:PROJECT_ID.svc.id.goog[timgcpsm-operator-system/timgcpsm-operator]"
```

Reinicia o pod do operador para aplicar a anotação, se já estiver em execução:

```bash
kubectl rollout restart deployment/timgcpsm-operator-controller -n timgcpsm-operator-system
```

---

## Fora do GKE (ou sem Workload Identity)

- **ADC** com credencial montada ou metadata (menos comum para este operador em prod): o processo precisa de uma identidade com `secretAccessor` equivalente.
- **Não** commits de JSON de conta de serviço no repositório; usa Secrets do cluster ou WI.

---

## Verificação rápida

- Logs do controller sem erros `permission denied` / `403` ao aceder ao GSM.
- Na GCP: **IAM** da GSA → papel `Secret Manager Secret Accessor` visível (direto ou herdado).
- `kubectl describe sa timgcpsm-operator -n timgcpsm-operator-system` → anotação `iam.gke.io/gcp-service-account` presente quando usas WI.

---

## Documentação Google

- [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
- [Secret Manager — controlo de acesso](https://cloud.google.com/secret-manager/docs/access-control)
- [Autenticação em GKE](https://cloud.google.com/kubernetes-engine/docs/concepts/workload-identity)
