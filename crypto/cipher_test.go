package crypto

import "testing"

func TestEncryptionNonce(t *testing.T) {
  cipherA := Encrypt("test text", "XkHBmh5Vvvk9OLGE8og9JbNvO3VWz2xJzjkRRfGyQ4Y=")

  cipherB := Encrypt("test text", "XkHBmh5Vvvk9OLGE8og9JbNvO3VWz2xJzjkRRfGyQ4Y=")

  if cipherA == cipherB {
    t.Errorf("The same text was cipher to the same output %s == %s", cipherA, cipherB)
  }
}

func TestEncode(t *testing.T) {
  encoding := encode([]byte{101, 102, 101, 101, 10, 43, 12})
  if encoding != "ZWZlZQorDA==" {
    t.Errorf("Encoding not working %s", encoding)
  }
}
