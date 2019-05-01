# Secrets
My secret management script

This project is not to replace commercial solutions. It is a simple solution to
grant some level of security for your personal files and store it in the cloud.

# Using

Generate a key

```
make gen
mkdir -p ~/.config/vault
mv vault.key ~/.config/vault/
chmod 0400 ~/.config/vault/vault.key
```

cipher a text:

```
./bin/vault \
  --k ~/.config/vault/vault.key \
  --tp plain.txt \
  --cp cipher.txt \
  --e
```

decipher a text

```
./bin/vault \
  --k ~/.config/vault/vault.key \
  --tp plain.txt \
  --cp cipher.txt \
  --d
```

# For future

[ ] Make all tests

[ ] Configure a cloud storage
[ ] Search for files in a cloud storage
[ ] Search keys in the default directory
[ ] Allow user key definition
[ ] Creates a http server
[ ] Creates REST API to cipher and decipher
[ ] Enforce the usage of ssl
[ ] Create a docker
[ ] Create a config file
