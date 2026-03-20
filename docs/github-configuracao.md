# Configuração do repositório no GitHub (segurança e OSS)

Este documento é para **mantenedores**: o que configurar na UI do GitHub para um projeto open source **apresentável e seguro**. Não precisas de “autenticação extra” no GitHub além da tua conta — o que importa é **como** a usas e as **definições do repositório e da organização**.

## Autenticação e contas

| O quê | Porquê |
| ----- | ------ |
| **2FA obrigatória** na tua conta (e, se tiveres org, para todos os membros com write) | Reduz takeover de conta e pushes maliciosos |
| **SSH keys** ou **GitHub CLI** com cuidado | Preferir chaves a passwords; revogar chaves antigas |
| **Fine-grained PATs** só quando necessário | Tokens clássicos com scope amplo são risco; PAT fino + expiração |

Não é preciso outro tipo de “login” especial para o repo em si: **branch protection**, **secrets** e **Actions** usam as permissões do GitHub normalmente.

## Definições recomendadas do repositório

### Geral (Settings → General)

- **Description, website, topics** — preencher (ex.: `kubernetes`, `gcp`, `secret-manager`, `operator`) para descoberta.
- **Features** — Issues e (opcional) Discussions se quiseres Q&A fora de issues.

### Segurança (Settings → Security)

- Ativar **Private vulnerability reporting** (recomendado) — alinha com [SECURITY.md](../SECURITY.md).
- **Dependabot alerts** — ativado por defeito na UI (alertas de CVE); **version updates** via `.github/dependabot.yml` estão **desativadas** neste repo.
- **Secret scanning** e **push protection** — disponíveis conforme o plano; úteis para evitar commits acidentais de tokens.

### Actions (Settings → Actions → General)

- **Actions permissions**: “Allow … actions and reusable workflows” — restringir a **marketplace e workflows da própria org** se quiseres mais controlo.
- **Fork pull requests**: por defeito, workflows de forks têm permissões mínimas; mantém política conservadora para repos públicos.

### Pacotes (GHCR)

- Garantir que a imagem `ghcr.io/.../timgcpsmoperator` tem **visibilidade** adequada (pública para pull open source).

## Proteção do ramo `main` (recomendado)

**Settings → Branches → Branch protection rule** para `main`:

- [ ] **Require a pull request before merging** (sem commits diretos ao `main`).
- [ ] **Require status checks to pass** — escolher o job de CI (ex.: `test` do workflow `CI`).
- [ ] **Require branches to be up to date before merging** (opcional, equipas maiores).
- [ ] **Require linear history** (opcional).
- [ ] **Do not allow bypassing the above settings** — pelo menos para não-admins.
- [ ] **Restrict who can push** — só tu / equipa core no `main`.

Isto **protege o repo** contra merges acidentais e mantém a história alinhada com CI verde.

## Secrets do repositório

- Usa **GitHub Secrets** só para o que os workflows precisam (ex.: `CODECOV_TOKEN` se passares a usar upload token no Codecov).
- **Nunca** commits com credenciais — o `GITHUB_TOKEN` nos workflows é injetado automaticamente e é de curta duração.

## O que já está no código deste repositório

- [LICENSE](../LICENSE) (Apache 2.0) e [NOTICE](../NOTICE)
- [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md), [CONTRIBUTING.md](../CONTRIBUTING.md), [SECURITY.md](../SECURITY.md)
- Dependabot version updates **desativadas** (sem `dependabot.yml`); alertas de segurança continuam configuráveis em **Settings → Security**
- Templates de [issue](../.github/ISSUE_TEMPLATE/) e [PR](../.github/pull_request_template.md)
- Workflows com permissões mínimas onde aplicável (`contents: read` no CI)

## Resumo

| Prioridade | Ação |
| ---------- | ---- |
| Alta | 2FA, branch protection no `main`, vulnerabilidades privadas |
| Média | Topics, alertas de dependências, revisar permissões de Actions |
| Baixa | Discussions, wiki, README social preview |

Se quiseres endurecer ainda mais: **revisões obrigatórias** (2 reviewers), **signed commits**, e políticas ao nível da **organização** (se aplicável).
