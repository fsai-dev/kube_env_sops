# kube_env_sops

Generate an encrypted Kubernetes Secret from a dotenv file using SOPS

## Installation

### Requires 
This package relies on `kubectl` and `sops` being installed and included in your $PATH

### Download and install from the shell

```shell
URL=https://github.com/e1-io/kube_env_sops/releases/download/v0.01/kube_env_sops-v0.01-$(uname -s)-$(uname -m).tar.gz
wget -qO- ${URL} | tar xvz - -C .
chmod +x kube_env_sops
sudo mv kube_env_sops /usr/local/bin
```

### Download and install from a release

[Releases · e1-io/kube_env_sops · GitHub](https://github.com/e1-io/kube_env_sops/releases)


## Usage

```shell
# Add the following to your .gitignore (recommended)
# files with the suffix -enc.y[a]ml are encrypted
# files with the suffix -dec.y[a]ml are encrypted
cat <<EOT > .gitignore
.env.local
*-dec.yml
*-dec.yaml
EOT

# Create a .env.local with variables
cat <<EOT > .env.local
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=admin
DATABASE_PASSWORD=password
DATABASE_NAME=postgres
EOT

# Generate an encrypted Kubernetes Secret from .env.local using SOPS
kube_env_sops
Successfully created the encrypted secret: .env-enc.yml

# If ./.env-enc.yml already exists, use the -force flag to overwrite
kube_env_sops -force

# Generate an encrypted Kubernetes Secret from .env.local using SOPS
kube_env_sops -save=false > .env-enc.yml
Outputs the contents of the encrpted manifest without saving

# Decrypt the encrypted file using sops and apply the manifest to kubernetes
SOPS_AGE_KEY_FILE=~/keys.txt sops --decrypt ./.env-enc.yml | kubectl apply -n default -f -
secret/environment-secrets created
```
