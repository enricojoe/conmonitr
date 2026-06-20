# Scripts

Utility scripts for local development and project setup.

---

## gen-secret.sh

Generates a cryptographically secure 40-character alphanumeric password and prints it as a `KEY=value` line ready to paste into a `.env` file.

**Usage**

```bash
./scripts/gen-secret.sh <ENV_VAR_NAME>
```

**Example**

```bash
./scripts/gen-secret.sh JWT_SECRET
# JWT_SECRET=fqFf7S0uFgDMEXOTe6AIRdp4Mtn96VW5BpLM59L1
```

Append directly to `.env`:

```bash
./scripts/gen-secret.sh JWT_SECRET >> .env
```
