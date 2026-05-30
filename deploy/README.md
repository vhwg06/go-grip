# Deploy Grip Backend To GCP Debian VM

This production deploy path targets the GCP VM `instance-20260526-070303`
running Debian 12. GitHub Actions builds the Docker image, pushes it to GHCR,
then SSHes into the VM and restarts Docker Compose from `/opt/go-grip`.

## Debian 12 VM Setup

Install Docker using Debian's official repository flow:

```sh
sudo apt-get update
sudo apt-get install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

Create the deploy directory and allow the SSH deploy user to run Docker. Replace
`quan` with the same username used by the SSH key and GitHub secret `VM_USER`.

```sh
sudo mkdir -p /opt/go-grip/deploy
sudo chown -R quan:quan /opt/go-grip
sudo usermod -aG docker quan
```

Log out and back in after adding the user to the `docker` group.

## GitHub Secrets

Required:

- `VM_HOST`
- `VM_USER`
- `VM_SSH_KEY`
- `GHCR_PAT` with `read:packages`
- `JWT_SECRET`
- `POSTGRES_PASSWORD`

Use a URL-safe `POSTGRES_PASSWORD` because the app currently receives
PostgreSQL through `PG_URL`, for example letters, numbers, `_`, and `-`.
Avoid reserved URL characters such as `@`, `/`, `:`, `?`, `#`, and spaces.

Optional:

- `VM_PORT`, defaults to `22`
- `VM_DEPLOY_DIR`, defaults to `/opt/go-grip`
- `ADMIN_USERS`
- `PAYMENT_MERCHANT_ID`
- `PAYMENT_SECRET_KEY`
- `PAYMENT_BASE_URL`
- `PAYMENT_NOTIFY_URL`
- `PAYMENT_RETURN_URL`

`VM_USER` must match the username used when creating or uploading the SSH key,
for example `quan`.

## GCP Firewall

Do not expose backend port `8080` publicly. Backend is bound to
`127.0.0.1:8080` on the VM and should be reached only through a reverse proxy
(frontend Nginx `/api` path or an edge proxy). PostgreSQL stays private inside
the Docker network and does not publish port `5432`.

## Smoke Test

After a deploy:

```sh
cd /opt/go-grip
docker compose -f deploy/docker-compose.prod.yml ps
curl -f http://127.0.0.1:8080/healthz
```

From outside the VM:

```sh
curl -f http://<VM_EXTERNAL_IP>/api/healthz
```
