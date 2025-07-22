package pool

import (
	"fmt"
	"sync"
)

type LazyReusePool struct {
	mtx    sync.Mutex
	head   *segment // nil when empty
	first  int
	last   int
	remain int
}

type segment struct {
	first int
	last  int
	next  *segment // nil when this segment is tail
}

type relativePos = int

const (
	withinThisSegment relativePos = iota
	adjacentToTheFront
	adjacentToTheBack
	before
	after
)

// NewLazyReusePool makes a LazyReusePool.
func NewLazyReusePool(first, last int) (*LazyReusePool, error) {
	if first > last {
		return nil, fmt.Errorf("make sure first(%d) <= last(%d)", first, last)
	}
	head := &segment{first, last, nil}
	return &LazyReusePool{
		head:   head,
		first:  first,
		last:   last,
		remain: last - first + 1,
	}, nil
}

func (p *LazyReusePool) Allocate() (res int, ok bool) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.head == nil {
		return 0, false
	}
	res = p.head.first
	p.head.first++
	if p.head.first > p.head.last {
		p.head = p.head.next
	}
	p.remain--
	return res, true
}

func (p *LazyReusePool) Use(value int) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if p.head == nil {
		return false
	}
	var prev *segment
	for cur := p.head; cur != nil; cur = cur.next {
		if cur.relativePosisionOf(value) == withinThisSegment {
			if cur.first == cur.last {
				if prev == nil {
					p.head = cur.next
				} else {
					prev.next = cur.next
				}
			} else {
				cur.split(value)
			}
			p.remain--
			return true
		}
		prev = cur
	}
	return false
}

func (p *LazyReusePool) Free(value int) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	// Ensure the value is within this pool
	if value < p.first || p.last < value {
		return false
	}

	// When this pool is empty
	if p.head == nil {
		p.head = newSingleSegment(value)
		p.remain++
		return true
	}

	// The value is returned into a segment after the head
	// for lazy reuse, excepted in the case of the value is
	// adjacent to the back of the head.

	switch p.head.relativePosisionOf(value) {
	case withinThisSegment:
		// duplecated free
		return false
	case adjacentToTheBack:
		// only in this case, returned into the head segment
		p.head.extendLast()
		p.remain++
		return true
	}

	// When there is only head segment
	if p.head.next == nil {
		// add a new segment
		p.head.next = newSingleSegment(value)
		p.remain++
		return true
	}

	prev := p.head
	cur := p.head.next
	for ; cur != nil; prev, cur = cur, cur.next {
		pos := cur.relativePosisionOf(value)
		switch pos {
		case before:
			// insert a segment
			temp := newSingleSegment(value)
			prev.next = temp
			temp.next = cur
			goto success
		case adjacentToTheFront:
			// extendFirst
			cur.first = value
			goto success
		case withinThisSegment:
			// duplecated free
			return false
		case adjacentToTheBack:
			cur.extendLast()
			goto success
		case after:
			if cur.next == nil {
				cur.next = newSingleSegment(value)
				goto success
			}
			// to next loop
		}
	}

success:
	p.remain++
	return true
}

func (p *LazyReusePool) Reserve(first, last int) error {
	if !p.Contains(first, last) {
		return fmt.Errorf("reserve range should in [%d, %d]", p.first, p.last)
	}

	for cur, prev := p.head, (*segment)(nil); cur != nil; cur = cur.next {
		switch {
		case cur.first >= first && cur.first <= last && cur.last > last:
			p.remain -= last - cur.first + 1
			cur.first = last + 1
		case cur.first < first && cur.last >= first && cur.last <= last:
			p.remain -= cur.last - first + 1
			cur.last = first - 1
		case cur.first < first && cur.last > last:
			cur.next = &segment{
				first: last + 1,
				last:  cur.last,
				next:  cur.next,
			}

			cur.last = first - 1
			p.remain -= last - first + 1

		// this segment in reserve range
		case cur.first >= first && cur.last <= last:
			p.remain -= cur.last - cur.first + 1
			if prev != nil {
				prev.next = cur.next
			} else {
				p.head = cur.next
			}
			// do not update prev
			continue
		}

		prev = cur
	}

	return nil
}

func (p *LazyReusePool) Contains(first, last int) bool {
	return first <= last && p.first <= first && p.last >= last
}

func (p *LazyReusePool) Min() int {
	return p.first
}

func (p *LazyReusePool) Max() int {
	return p.last
}

func (p *LazyReusePool) Remain() int {
	return p.remain
}

func (p *LazyReusePool) Total() int {
	return p.last - p.first + 1
}

func (p *LazyReusePool) GetHead() *segment {
	return p.head
}

func newSingleSegment(num int) *segment {
	return &segment{num, num, nil}
}

func (s *segment) relativePosisionOf(value int) relativePos {
	switch {
	case value < s.first-1:
		return before
	case value == s.first-1:
		return adjacentToTheFront
	case s.first <= value && value <= s.last:
		return withinThisSegment
	case value == s.last+1:
		return adjacentToTheBack
	default:
		return after
	}
}

func (s *segment) split(use int) bool {
	if use < s.first || use > s.last {
		return false
	}

	switch use {
	case s.first:
		s.first += 1
	case s.last:
		s.last -= 1
	default:
		next := &segment{
			first: use + 1,
			last:  s.last,
			next:  s.next,
		}

		s.next = next
		s.last = use - 1
	}

	return true
}

func (s *segment) extendLast() *segment {
	s.last++
	if s.next != nil && s.last+1 == s.next.first {
		// concatenate
		s.last = s.next.last
		s.next = s.next.next
	}
	return s
}

func (s *segment) First() int {
	return s.first
}

func (s *segment) Last() int {
	return s.last
}

func (s *segment) Next() *segment {
	return s.next
}

func (p1 *LazyReusePool) IsJoint(p2 *LazyReusePool) bool {
	if p2.last < p1.first || p1.last < p2.first {
		return false
	}
	return true
}

func (p *LazyReusePool) Dump() [][]int {
	var dumpedSegList [][]int
	curSeg := p.head
	for curSeg != nil {
		dumpedSeg := []int{curSeg.first, curSeg.last}
		dumpedSegList = append(dumpedSegList, dumpedSeg)
		curSeg = curSeg.next
	}
	return dumpedSegList
}
