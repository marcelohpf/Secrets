CC=go
PASS=meussegredos
BINDIR=./bin

vault: $(BINDIR)
	$(CC) build -o $(BINDIR)/vault

$(BINDIR):
	mkdir -p $(BINDIR)

run: vault cipher decipher
cipher: vault
	$(BINDIR)/vault seal --key examples/vault.key --text-path examples/plain.txt --cipher-path examples/cipher.txt

decipher: vault
	$(BINDIR)/vault unseal --key  examples/vault.key --cipher-path examples/cipher.txt


gen: vault
	$(BINDIR)/vault keygen --key-path ./examples/ --key vault.key

.PHONY: clean vault

clean:
	rm -rf $(BINDIR)/vault
