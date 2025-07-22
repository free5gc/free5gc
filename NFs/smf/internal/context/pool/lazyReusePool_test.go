package pool_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/free5gc/smf/internal/context/pool"
)

func TestNewLazyReusePool(t *testing.T) {
	// OK : first < last
	lrp, err := pool.NewLazyReusePool(10, 100)
	assert.NoError(t, err)
	assert.NotEmpty(t, lrp)
	assert.Equal(t, 91, lrp.Total())
	assert.Equal(t, 91, lrp.Remain())

	// OK : first == last
	lrp, err = pool.NewLazyReusePool(100, 100)
	assert.NoError(t, err)
	assert.NotEmpty(t, lrp)
	assert.Equal(t, 1, lrp.Total())
	assert.Equal(t, 1, lrp.Remain())

	// NG : first > last
	lrp, err = pool.NewLazyReusePool(10, 0)
	assert.Empty(t, lrp)
	assert.Error(t, err)
}

func TestLazyReusePool_SingleSegment(t *testing.T) {
	// Allocation OK
	p, err := pool.NewLazyReusePool(1, 2)
	assert.NoError(t, err)
	a, ok := p.Allocate()
	assert.Equal(t, 1, a)
	assert.True(t, ok)
	assert.Equal(t, 1, p.Remain())
	a, ok = p.Allocate()
	assert.Equal(t, a, 2)
	assert.True(t, ok)
	assert.Equal(t, 0, p.Remain())

	// exhauted
	_, ok = p.Allocate()
	assert.False(t, ok)
	assert.Equal(t, 0, p.Remain())

	// free 1
	ok = p.Free(1)
	assert.True(t, ok)

	assert.Equal(t, 1, p.GetHead().First())
	assert.Equal(t, 1, p.GetHead().Last())
	assert.Empty(t, p.GetHead().Next())
	assert.Equal(t, 1, p.Remain())

	// out of range
	ok = p.Free(0)
	assert.False(t, ok)

	// duplecated free
	ok = p.Free(1)
	assert.False(t, ok)

	// free 2
	ok = p.Free(2)
	assert.True(t, ok)
	assert.Equal(t, 2, p.Remain())

	// reuse
	a, ok = p.Allocate()
	assert.Equal(t, 1, a)
	assert.True(t, ok)
	assert.Equal(t, 1, p.Remain())
	assert.Equal(t, 2, p.Total())

	ok = p.Free(1)
	assert.True(t, ok)
	assert.Equal(t, 2, p.Remain())

	ok = p.Use(2)
	assert.True(t, ok)
	assert.Equal(t, 1, p.Remain())

	ok = p.Use(1)
	assert.True(t, ok)
	assert.Equal(t, 0, p.Remain())
	assert.Equal(t, 2, p.Total())

	// try use from empty pool
	ok = p.Use(1)
	assert.False(t, ok)
	assert.Equal(t, 0, p.Remain())
	assert.Equal(t, 2, p.Total())

	ok = p.Free(1)
	assert.True(t, ok)
	assert.Equal(t, 1, p.Remain())
	assert.Equal(t, 2, p.Total())

	// try use from assigned value
	ok = p.Use(2)
	assert.False(t, ok)
	assert.Equal(t, 1, p.Remain())
	assert.Equal(t, 2, p.Total())

	ok = p.Free(2)
	assert.True(t, ok)
	assert.Equal(t, 2, p.Remain())
	assert.Equal(t, 2, p.Total())

	// split from s.last
	ok = p.Use(2)
	assert.True(t, ok)
	assert.Equal(t, 1, p.Remain())
	assert.Equal(t, 2, p.Total())
}

func TestLazyReusePool_ManySegment(t *testing.T) {
	p, err := pool.NewLazyReusePool(1, 100)
	assert.NoError(t, err)
	assert.Equal(t, 100, p.Remain())

	// -> 1-100

	for i := 0; i < 99; i++ {
		_, ok := p.Allocate()
		assert.True(t, ok)
	}

	// -> 100-100
	assert.Equal(t, 1, p.Remain())

	p.Free(3) // -> 100-100 -> 3-3
	assert.Equal(t, 3, p.GetHead().Next().First())
	assert.Equal(t, 3, p.GetHead().Next().Last())

	p.Free(6) // -> 100-100 -> 3-3 -> 6-6
	assert.Equal(t, 6, p.GetHead().Next().Next().First())
	assert.Equal(t, 6, p.GetHead().Next().Next().Last())
	assert.Equal(t, 3, p.Remain())

	// adjacent to the front
	p.Free(2) // -> 100-100 -> 2-3 -> 6-6
	assert.Equal(t, 2, p.GetHead().Next().First())
	assert.Equal(t, 3, p.GetHead().Next().Last())
	assert.Equal(t, 4, p.Remain())

	// duplicate
	ok := p.Free(3)
	assert.False(t, ok)

	// adjacent to the back
	p.Free(4) // -> 100-100 -> 2-4 -> 6-6
	assert.Equal(t, 2, p.GetHead().Next().First())
	assert.Equal(t, 4, p.GetHead().Next().Last())
	assert.Equal(t, 5, p.Remain())

	// 3rd segment
	p.Free(7) // -> 100-100 -> 2-4 -> 6-7
	assert.Equal(t, 6, p.GetHead().Next().Next().First())
	assert.Equal(t, 7, p.GetHead().Next().Next().Last())
	assert.Equal(t, 6, p.Remain())

	// concatenate
	p.Free(5) // -> 100-100 -> 2-7
	assert.Equal(t, 2, p.GetHead().Next().First())
	assert.Equal(t, 7, p.GetHead().Next().Last())
	assert.Empty(t, p.GetHead().Next().Next())
	assert.Equal(t, 7, p.Remain())

	// new head
	a, ok := p.Allocate() // -> 2-7
	assert.Equal(t, a, 100)
	assert.True(t, ok)
	assert.Equal(t, 2, p.GetHead().First())
	assert.Equal(t, 7, p.GetHead().Last())
	assert.Equal(t, 6, p.Remain())

	// reuse
	a, ok = p.Allocate() // -> 3-7
	assert.Equal(t, a, 2)
	assert.True(t, ok)
	assert.Equal(t, 3, p.GetHead().First())
	assert.Equal(t, 7, p.GetHead().Last())
	assert.Empty(t, p.GetHead().Next())
	assert.Equal(t, 5, p.Remain())

	// return to head
	p.Free(100) // -> 3-7 -> 100-100
	p.Free(8)   // -> 3-8 -> 100-100
	assert.Equal(t, 3, p.GetHead().First())
	assert.Equal(t, 8, p.GetHead().Last())
	assert.Equal(t, 100, p.GetHead().Next().First())
	assert.Equal(t, 100, p.GetHead().Next().Last())
	assert.Equal(t, 7, p.Remain())

	// return into between head and 2nd
	p.Free(10) // -> 3-8 -> 10-10 -> 100-100
	assert.Equal(t, 3, p.GetHead().First())
	assert.Equal(t, 8, p.GetHead().Last())
	assert.Equal(t, 10, p.GetHead().Next().First())
	assert.Equal(t, 10, p.GetHead().Next().Last())
	assert.Equal(t, 100, p.GetHead().Next().Next().First())
	assert.Equal(t, 100, p.GetHead().Next().Next().Last())
	assert.Equal(t, 8, p.Remain())

	// concatenate head and 2nd
	p.Free(9) // -> 3-10 -> 100-100
	assert.Equal(t, 3, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Equal(t, 100, p.GetHead().Next().First())
	assert.Equal(t, 100, p.GetHead().Next().Last())
	assert.Empty(t, p.GetHead().Next().Next())
	assert.Equal(t, 9, p.Remain())

	ok = p.Use(100) // 3-10
	assert.True(t, ok)
	assert.Equal(t, 3, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Nil(t, p.GetHead().Next())
	assert.Equal(t, 8, p.Remain())

	ok = p.Use(3) // 4-10
	assert.True(t, ok)
	assert.Equal(t, 4, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Nil(t, p.GetHead().Next())
	assert.Equal(t, 7, p.Remain())

	ok = p.Use(7) // 4-6 -> 8-10
	assert.True(t, ok)
	assert.Equal(t, 4, p.GetHead().First())
	assert.Equal(t, 6, p.GetHead().Last())
	assert.Equal(t, 8, p.GetHead().Next().First())
	assert.Equal(t, 10, p.GetHead().Next().Last())
	assert.Equal(t, 6, p.Remain())

	p.Free(7) // 4-10
	assert.Equal(t, 4, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Nil(t, p.GetHead().Next())

	p.Free(3) // 4-10 -> 3
	assert.Equal(t, 4, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Equal(t, 3, p.GetHead().Next().First())
	assert.Equal(t, 3, p.GetHead().Next().Last())

	p.Free(100) // 4-10 -> 3 -> 100
	assert.Equal(t, 4, p.GetHead().First())
	assert.Equal(t, 10, p.GetHead().Last())
	assert.Equal(t, 3, p.GetHead().Next().First())
	assert.Equal(t, 3, p.GetHead().Next().Last())
	assert.Equal(t, 100, p.GetHead().Next().Next().First())
	assert.Equal(t, 100, p.GetHead().Next().Next().Last())
}

func TestLazyReusePool_ReserveSection(t *testing.T) {
	p, err := pool.NewLazyReusePool(1, 100)
	require.NoError(t, err)
	require.Equal(t, 100, p.Remain())

	err = p.Reserve(30, 69)
	require.NoError(t, err)
	require.Equal(t, 60, p.Remain())

	var allocated []int
	for {
		a, ok := p.Allocate()
		if ok {
			allocated = append(allocated, a)
		} else {
			break
		}
	}
	var expected []int
	for i := 1; i <= 29; i++ {
		expected = append(expected, i)
	}
	for i := 70; i <= 100; i++ {
		expected = append(expected, i)
	}

	require.Equal(t, expected, allocated)
}

func TestLazyReusePool_ReserveSection2(t *testing.T) {
	p, err := pool.NewLazyReusePool(10, 100)
	require.NoError(t, err)
	assert.Equal(t, (100 - 10 + 1), p.Remain())

	assert.Equal(t, p.GetHead().First(), 10)
	assert.Equal(t, p.GetHead().Last(), 100)

	// try reserve outside range
	err = p.Reserve(0, 5)
	assert.Error(t, err)

	// reserve entries on head
	err = p.Reserve(10, 20)
	require.NoError(t, err)
	assert.Equal(t, (100 - 21 + 1), p.Remain())

	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 100)

	// reserve entries on tail
	err = p.Reserve(90, 100)
	require.NoError(t, err)
	assert.Equal(t, (89 - 21 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 89)

	// reserve entries on center
	err = p.Reserve(40, 50)
	require.NoError(t, err)
	assert.Equal(t, (39 - 21 + 1 + 89 - 51 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 39)
	assert.Equal(t, p.GetHead().Next().First(), 51)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	// try reserve range was already reserved
	err = p.Reserve(10, 20)
	require.NoError(t, err)
	assert.Equal(t, (39 - 21 + 1 + 89 - 51 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 39)
	assert.Equal(t, p.GetHead().Next().First(), 51)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	err = p.Reserve(40, 50)
	require.NoError(t, err)
	assert.Equal(t, (39 - 21 + 1 + 89 - 51 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 39)
	assert.Equal(t, p.GetHead().Next().First(), 51)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	err = p.Reserve(90, 100)
	require.NoError(t, err)
	assert.Equal(t, (39 - 21 + 1 + 89 - 51 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 39)
	assert.Equal(t, p.GetHead().Next().First(), 51)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	// reserve range includes reserved and non-reserved addresses
	err = p.Reserve(36, 55)
	require.NoError(t, err)
	assert.Equal(t, (35 - 21 + 1 + 89 - 56 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 21)
	assert.Equal(t, p.GetHead().Last(), 35)
	assert.Equal(t, p.GetHead().Next().First(), 56)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	// remove entire segment
	err = p.Reserve(21, 35)
	require.NoError(t, err)
	assert.Equal(t, (89 - 56 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 56)
	assert.Equal(t, p.GetHead().Last(), 89)

	// generate 3 segments
	err = p.Reserve(70, 75)
	require.NoError(t, err)
	assert.Equal(t, (69 - 56 + 1 + 89 - 76 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 56)
	assert.Equal(t, p.GetHead().Last(), 69)
	assert.Equal(t, p.GetHead().Next().First(), 76)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	err = p.Reserve(60, 65)
	require.NoError(t, err)
	assert.Equal(t, (59 - 56 + 1 + 69 - 66 + 1 + 89 - 76 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 56)
	assert.Equal(t, p.GetHead().Last(), 59)
	assert.Equal(t, p.GetHead().Next().First(), 66)
	assert.Equal(t, p.GetHead().Next().Last(), 69)
	assert.Equal(t, p.GetHead().Next().Next().First(), 76)
	assert.Equal(t, p.GetHead().Next().Next().Last(), 89)

	// remove center segment
	err = p.Reserve(60, 75)
	require.NoError(t, err)
	assert.Equal(t, (59 - 56 + 1 + 89 - 76 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 56)
	assert.Equal(t, p.GetHead().Last(), 59)
	assert.Equal(t, p.GetHead().Next().First(), 76)
	assert.Equal(t, p.GetHead().Next().Last(), 89)

	// remove tail segment
	err = p.Reserve(70, 90)
	require.NoError(t, err)
	assert.Equal(t, (59 - 56 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 56)
	assert.Equal(t, p.GetHead().Last(), 59)

	// remove last segment
	err = p.Reserve(50, 60)
	require.NoError(t, err)
	assert.Equal(t, 0, p.Remain())
	assert.Nil(t, p.GetHead())
}

func TestLazyReusePool_ReserveSection3(t *testing.T) {
	p, err := pool.NewLazyReusePool(10, 99)
	require.NoError(t, err)
	assert.Equal(t, (99 - 10 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 10)
	assert.Equal(t, p.GetHead().Last(), 99)

	// generate 4 segments
	err = p.Reserve(20, 29)
	require.NoError(t, err)
	err = p.Reserve(40, 49)
	require.NoError(t, err)
	err = p.Reserve(60, 69)
	require.NoError(t, err)
	require.Equal(t, (19 - 10 + 1 + 39 - 30 + 1 + 59 - 50 + 1 + 99 - 70 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 10)
	assert.Equal(t, p.GetHead().Last(), 19)
	assert.Equal(t, p.GetHead().Next().First(), 30)
	assert.Equal(t, p.GetHead().Next().Last(), 39)
	assert.Equal(t, p.GetHead().Next().Next().First(), 50)
	assert.Equal(t, p.GetHead().Next().Next().Last(), 59)
	assert.Equal(t, p.GetHead().Next().Next().Next().First(), 70)
	assert.Equal(t, p.GetHead().Next().Next().Next().Last(), 99)

	// remove two segments
	err = p.Reserve(30, 59)
	require.NoError(t, err)
	require.Equal(t, (19 - 10 + 1 + 99 - 70 + 1), p.Remain())
	assert.Equal(t, p.GetHead().First(), 10)
	assert.Equal(t, p.GetHead().Last(), 19)
	assert.Equal(t, p.GetHead().Next().First(), 70)
	assert.Equal(t, p.GetHead().Next().Last(), 99)
}

func TestLazyReusePool_ManyGoroutine(t *testing.T) {
	p, err := pool.NewLazyReusePool(101, 1000)
	assert.NoError(t, err)
	assert.Equal(t, 900, p.Remain())
	ch := make(chan int)

	numOfThreads := 400

	for i := 0; i < numOfThreads; i++ {
		// Allocate 2 times and Free 1 time
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// Program will be exited.
					fmt.Fprintf(os.Stderr, "panic: %v\n%s", p, string(debug.Stack()))
					os.Exit(1)
				}
			}()

			a1, ok := p.Allocate()
			assert.True(t, ok)
			ch <- a1

			time.Sleep(10 * time.Millisecond)

			ok = p.Free(a1)
			assert.True(t, ok)

			time.Sleep(10 * time.Millisecond)

			a2, ok := p.Allocate()
			assert.True(t, ok)
			ch <- a2
		}()
	}
	// collect allocated values
	allocated := make([]int, 0, numOfThreads*2)
	for i := 0; i < numOfThreads*2; i++ {
		allocated = append(allocated, <-ch)
	}
	sort.Ints(allocated)

	expected := make([]int, numOfThreads*2)
	for i := 0; i < numOfThreads*2; i++ {
		expected[i] = p.Min() + i
	}
	assert.Equal(t, expected, allocated)
	assert.Equal(t, 900-numOfThreads, p.Remain())

	a, ok := p.Allocate()
	assert.Equal(t, p.Min()+numOfThreads*2, a)
	assert.True(t, ok)
	assert.Equal(t, 900-numOfThreads-1, p.Remain())
}
