package crypto

import (
  "testing"
)

func TestDifferentGenerator(t *testing.T) {
  data := generate(32)
  data2 := generate(32)
  equal := 0

  for i := 0; i < 32; i++ { 
    if data[i] == data2[i] {
      equal++
    }
  }

  if equal == 32 {
    t.Errorf("Data generated is equal")
  }
}

func TestSizeKeyGenerator(t *testing.T) {

  var sizes []int = []int{16, 32}

  for _, size := range sizes {
    data := generate(size)
    if len(data) != size {
      t.Errorf("Data generated is different of the expected %d!=%d", size, len(data))
    }
  }
}
