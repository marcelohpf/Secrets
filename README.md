# Secrets
My secret management script

This project is not to replace commercial solutions. It is a simple solution to
grant some level of security for your personal files and store it in the cloud.

# Using

Install Go and configure the `$GOPATH` before continue.

Building

```bash
git clone https://github.com/marcelohpf/secrets.git $GOPATH/github.com/marcelohpf/secrets
cd $GOPATH/github.com/marcelohpf/secrets
make build
make install
make clean
mkdir -p ~/vault
```

Generate a key

```bash
vault keygen --key vault.key
sudo chmod 0400 ~/vault/vault.key
```

Seal a secret

```bash
vault box seal --key vault.key --item secrets --in ./plain_in.txt
```

Unseal a secret

```bash
vault box unseal --key vault.key --item secrets --out ./plain_out.txt
```

Verify the result with

```bash
diff ./plain_in.txt ./plain_out.txt
```

# For future

- [ ] Make all tests

- [ ] Allow secret replication
- [ ] Add metadata
- [x] Seal/Unseal boxes in a cloud storage
- [x] Configure a cloud storage
- [ ] Allow user key definition
- [x] Creates a HTTP server
- [ ] Creates REST API to seal/unseal
- [ ] Create a docker for the HTTP server
- [ ] Create a config file
- [x] Search keys in the default directory
- [x] Use a default directory for itens of a box

# Google Drive Integration


