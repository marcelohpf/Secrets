CC=go
PASS=meussegredos
BINDIR=./bin

vault: $(BINDIR)
	$(CC) build -o $(BINDIR)/vault

$(BINDIR):
	mkdir -p $(BINDIR)

test: gen seal unseal
seal: vault
	$(BINDIR)/vault box seal --debug --key-path ./examples/ --key vault.key --box-path examples --item little-secrets --in ./examples/plain.txt

unseal: vault
	$(BINDIR)/vault box unseal --debug --key-path  ./examples/ --key vault.key --box-path examples --item little-secrets --out ./examples/plain.txt


gen: vault
	$(BINDIR)/vault genkey --debug --key-path ./examples/ --key-name vault.key

server: vault
	./bin/vault server --debug --key-path ./examples/ --box-path examples
.PHONY: clean vault

clean:
	rm -rf $(BINDIR)/ examples/little-secrets examples/vault.key
