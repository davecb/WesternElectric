package movingAverage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMovingAverage(t *testing.T) {
	assert := assert.New(t)

	add := New(5)
	a, b := add(1)
	t.Logf("(1                  ) / 5 = %g", a)
	assert.Equal(0.2, a)
	assert.Equal(0.447213595499958, b)

	a, b = add(2)
	t.Logf("(1+2                ) / 5 = %g", a)
	assert.Equal(0.6, a)
	assert.Equal(0.8944271909999159, b)

	a, b = add(3)
	t.Logf("(1+2+3              ) / 5 = %g", a)
	assert.Equal(1.2, a)
	assert.Equal(1.3038404810405297, b)

	a, b = add(4)
	t.Logf("(1+2+3+4            ) / 5 = %g", a)
	assert.Equal(2.0, a)
	assert.Equal(1.5811388300841898, b)

	a, b = add(5)
	t.Logf("(1+2+3+4+5          ) / 5 = %g", a)
	assert.Equal(3.0, a)
	assert.Equal(1.5811388300841898, b)

	a, b = add(9)
	t.Logf("(  2+3+4+5+9        ) / 5 = %g", a)
	assert.Equal(4.6, a)
	assert.Equal(2.701851217221259, b)

	a, b = add(3)
	t.Logf("(    3+4+5+9+3      ) / 5 = %g", a)
	assert.Equal(4.8, a)
	assert.Equal(2.4899799195977463, b)

	a, b = add(0)
	t.Logf("(      4+5+9+3+0    ) / 5 = %g", a)
	assert.Equal(4.2, a)
	assert.Equal(3.271085446759225, b)

	a, b = add(-9)
	t.Logf("(        5+9+3+0-9  ) / 5 = %g", a)
	assert.Equal(1.6, a)
	assert.Equal(6.767569726275452, b)

	a, b = add(-8)
	t.Logf("(          9+3+0-9-8) / 5 = %g", a)
	assert.Equal(-1.0, a)
	assert.Equal(7.582875444051551, b)

}

func ExampleMovingAverage() {

	add := New(5)
	a, _ := add(1)
	fmt.Println("(1                  ) / 5 =", a)
	a, _ = add(2)
	fmt.Println("(1+2                ) / 5 =", a)
	a, _ = add(3)
	fmt.Println("(1+2+3              ) / 5 =", a)
	a, _ = add(4)
	fmt.Println("(1+2+3+4            ) / 5 =", a)
	a, _ = add(5)
	fmt.Println("(1+2+3+4+5          ) / 5 =", a)
	a, _ = add(9)
	fmt.Println("(  2+3+4+5+9        ) / 5 =", a)
	a, _ = add(3)
	fmt.Println("(    3+4+5+9+3      ) / 5 =", a)
	a, _ = add(0)
	fmt.Println("(      4+5+9+3+0    ) / 5 =", a)
	a, _ = add(-9)
	fmt.Println("(        5+9+3+0-9  ) / 5 =", a)
	a, _ = add(-8)
	fmt.Println("(          9+3+0-9-8) / 5 =", a)
	// Output:
	//(1                  ) / 5 = 0.2
	//(1+2                ) / 5 = 0.6
	//(1+2+3              ) / 5 = 1.2
	//(1+2+3+4            ) / 5 = 2
	//(1+2+3+4+5          ) / 5 = 3
	//(  2+3+4+5+9        ) / 5 = 4.6
	//(    3+4+5+9+3      ) / 5 = 4.8
	//(      4+5+9+3+0    ) / 5 = 4.2
	//(        5+9+3+0-9  ) / 5 = 1.6
	//(          9+3+0-9-8) / 5 = -1
}

// Expected:
//(1+2+3+4+5          ) / 5 = 3
//(  2+3+4+5+9        ) / 5 = 4.6
//(    3+4+5+9+3      ) / 5 = 4.8
//(      4+5+9+3+0    ) / 5 = 4.2
//(        5+9+3+0-9  ) / 5 = 1.6
//(          9+3+0-9-8) / 5 = -1
