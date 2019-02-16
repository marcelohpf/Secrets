CC=go
PASS=meussegredos

vault:
	$(CC) build

run: vault cipher decipher
cipher: vault
	./vault --k vault.key --tp plain.txt --cp cipher.txt --e

decipher: vault
	./vault --k  vault.key --cp cipher.txt --d


gen: vault
	./vault --gk --k vault.key

.PHONY: clean vault

clean:
	rm vault
