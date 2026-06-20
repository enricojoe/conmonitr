# Scripts

Utility scripts for local development and project setup.

---

## gen-password.sh

Generates a random 40-character password **and its bcrypt hash**, outputting both lines ready to paste into `.env`. Use this for `CONMONITR_PASSWORD` / `CONMONITR_PASSWORD_HASH` — the backend validates login using bcrypt, so the hash must be pre-computed here.

**Usage**

```bash
./scripts/gen-password.sh
```

**Example output**

```
CONMONITR_PASSWORD=Qk6t9nuwgEpkosQgTtWYNcAORxPV5r2GoLtSA69T
CONMONITR_PASSWORD_HASH=$2a$10$ng0gG5KN6eny0BVvHxMgpOmYm0s1i4fUR33atoMdj7pFOqE.QTWLK
```

Paste both lines into `.env`. Use the plain `CONMONITR_PASSWORD` value when logging in; never put it in the env file (keep it somewhere safe). Only `CONMONITR_PASSWORD_HASH` is needed by the server.

Append directly to `.env`:

```bash
./scripts/gen-password.sh >> .env
```

---

## gen-secret.sh

Generates a cryptographically secure 40-character alphanumeric string as a `KEY=value` line. Use this for non-password secrets like `CONMONITR_JWT_SECRET`.

**Usage**

```bash
./scripts/gen-secret.sh <ENV_VAR_NAME>
```

**Example**

```bash
./scripts/gen-secret.sh CONMONITR_JWT_SECRET
# CONMONITR_JWT_SECRET=fqFf7S0uFgDMEXOTe6AIRdp4Mtn96VW5BpLM59L1
```

Append directly to `.env`:

```bash
./scripts/gen-secret.sh CONMONITR_JWT_SECRET >> .env
```
