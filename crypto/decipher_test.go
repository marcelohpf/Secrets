package crypto

import "testing"

func TestDecipher(t *testing.T) {
  plainText := Decrypt("WC0gnAzc689raPovnXCuq+yLS/AfoBWvIrckerefC2yLhFHs0A==", "XkHBmh5Vvvk9OLGE8og9JbNvO3VWz2xJzjkRRfGyQ4Y=")
  if "test text" != plainText {
    t.Errorf("The decryption fails test text != %s", plainText)
  }
}
