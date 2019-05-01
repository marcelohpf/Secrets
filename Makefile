CC=go
PASS=meussegredos
BINDIR=./bin

vault: $(BINDIR)
	$(CC) build -o $(BINDIR)/vault

$(BINDIR):
	mkdir -p $(BINDIR)

test: gen seal unseal
seal: vault
	$(BINDIR)/vault box seal --key-path ./examples/ --key vault.key --box-path ./examples --item little-secrets --in ./examples/plain.txt

unseal: vault
	$(BINDIR)/vault box unseal --key-path  ./examples/ --key vault.key --box-path ./examples/ --item little-secrets --out ./examples/plain.txt


gen: vault
	$(BINDIR)/vault keygen --key-path ./examples/ --key vault.key

.PHONY: clean vault

clean:
	rm -rf $(BINDIR)/vault
