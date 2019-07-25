package system

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/aergoio/aergo/state"
	"github.com/aergoio/aergo/types"
)

const (
	vprMax = uint32(50000)
	// vprBucketsMax must be smaller than 256. A bigger number is regarded as
	// 256.
	vprBucketsMax = 71
)

var (
	vprKeyPrefix = []byte("VotingPowerBucket/")
	zeroValue    = &big.Int{}

	votingPowerRank = newVpr()
)

type votingPower struct {
	addr  types.AccountID
	power *big.Int
}

func newVotingPower(addr types.AccountID, pow *big.Int) *votingPower {
	return &votingPower{addr: addr, power: pow}
}

func (vp *votingPower) getAddr() types.AccountID {
	return vp.addr
}

func (vp *votingPower) addrBytes() []byte {
	return vp.addr[:]
}

func (vp *votingPower) setAddr(addr types.AccountID) {
	vp.addr = addr
}

func (vp *votingPower) getPower() *big.Int {
	return vp.power
}

func (vp *votingPower) setPower(pow *big.Int) {
	vp.power = pow
}

func (vp *votingPower) cmp(pow *big.Int) int {
	return vp.getPower().Cmp(pow)
}

func (vp *votingPower) lt(rhs *votingPower) bool {
	return vp.cmp(rhs.getPower()) < 0
}

func (vp *votingPower) le(rhs *votingPower) bool {
	return vp.lt(rhs) || vp.cmp(rhs.getPower()) == 0
}

func (vp *votingPower) marshal() []byte {
	var buf bytes.Buffer

	buf.Write(vp.addrBytes())
	buf.Write(vp.getPower().Bytes())

	hdr := make([]byte, 4)
	binary.LittleEndian.PutUint32(hdr, uint32(buf.Len()))

	return append(hdr, buf.Bytes()...)
}

func (vp *votingPower) unmarshal(b []byte) uint32 {
	var n uint32

	r := bytes.NewReader(b)
	binary.Read(r, binary.LittleEndian, &n)

	vp.setAddr(types.AccountID(types.ToHashID(b[4:36])))
	if int(4+n) < len(b) {
		vp.setPower(new(big.Int).SetBytes(b[36 : 4+n]))
	} else {
		vp.setPower(new(big.Int).SetBytes(b[36:]))
	}

	return 4 + n
}

type vprStore struct {
	buckets map[uint8]*list.List
	cmp     func(lhs *votingPower, rhs *votingPower) int
	lowest  *votingPower
}

func newVprStore(bucketsMax uint32) *vprStore {
	return &vprStore{
		buckets: make(map[uint8]*list.List, bucketsMax),
		cmp: func(lhs *votingPower, rhs *votingPower) int {
			return bytes.Compare(lhs.addrBytes(), rhs.addrBytes())
		},
	}
}

func (b *vprStore) getLowest() *votingPower {
	return b.lowest
}

func (b *vprStore) update(vp *votingPower) (idx uint8) {
	var (
		addr = vp.getAddr()
		pow  = vp.getPower()
	)

	idx = getBucketIdx(addr)

	var (
		bu    *list.List
		exist bool
	)

	if bu, exist = b.buckets[idx]; !exist {
		bu = list.New()
		b.buckets[idx] = bu
	}

	v := remove(bu, addr)
	if v != nil {
		v.setPower(pow)
	} else {
		v = vp
	}

	orderedListAdd(bu, v,
		func(e *list.Element) bool {
			ind := b.cmp(toVotingPower(e), v)
			return ind <= 0
		},
	)

	b.updateLowest(v)

	return
}

type predicate func(e *list.Element) bool

func search(l *list.List, match func(e *list.Element) bool) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if match(e) {
			return e
		}
	}
	return nil
}

func orderedListElemAdd(l *list.List, e *list.Element, match predicate) {
	if m := search(l, match); m != nil {
		l.MoveBefore(e, m)
	} else {
		l.MoveToBack(e)
	}
}

func orderedListAdd(l *list.List, e *votingPower, match predicate) *list.Element {
	var voter *list.Element

	if m := search(l, match); m != nil {
		voter = l.InsertBefore(e, m)
	} else {
		voter = l.PushBack(e)
	}

	return voter
}

func (b *vprStore) updateLowest(v *votingPower) {
	if b.lowest == nil {
		b.lowest = v
	} else if v.lt(b.lowest) {
		b.lowest = v
	}
}

func (b *vprStore) addTail(i uint8, vp *votingPower) {
	var l *list.List
	if l = b.buckets[i]; l == nil {
		l = list.New()
		b.buckets[i] = l
	}
	l.PushBack(vp)
	b.updateLowest(vp)
}

type dataSetter interface {
	SetData(key, value []byte) error
}

type dataGetter interface {
	GetData(key []byte) ([]byte, error)
}

func (b *vprStore) write(s dataSetter, i uint8) error {
	var buf bytes.Buffer

	l := b.buckets[i]
	for e := l.Front(); e != nil; e = e.Next() {
		buf.Write(toVotingPower(e).marshal())
	}

	return s.SetData(vprKey(i), buf.Bytes())
}

func (b *vprStore) read(s dataGetter, i uint8) ([]*votingPower, error) {
	buf, err := s.GetData(vprKey(i))
	if err != nil {
		return nil, err
	}
	vps := make([]*votingPower, 0, 10)
	for off := 0; off < len(buf); {
		vp := &votingPower{}
		off += int(vp.unmarshal(buf[off:]))
		vps = append(vps, vp)
	}
	return vps, nil
}

func remove(bu *list.List, addr types.AccountID) *votingPower {
	for e := bu.Front(); e != nil; e = e.Next() {
		if v := toVotingPower(e); addr == v.getAddr() {
			return bu.Remove(e).(*votingPower)
		}
	}
	return nil
}

func getBucketIdx(addr types.AccountID) uint8 {
	return uint8(addr[0]) % vprBucketsMax
}

func (b *vprStore) getBucket(addr types.AccountID) *list.List {
	if b, exist := b.buckets[getBucketIdx(addr)]; exist {
		return b
	}
	return nil
}

type topVoters struct {
	list   *list.List
	cmp    func(lhs *big.Int, rhs *big.Int) bool
	max    uint32
	powers map[types.AccountID]*list.Element
}

func newTopVoters(max uint32) *topVoters {
	return &topVoters{
		list: list.New(),
		cmp: func(curr *big.Int, e *big.Int) bool {
			return curr.Cmp(e) >= 0
		},
		max:    max,
		powers: make(map[types.AccountID]*list.Element),
	}
}

func (tv *topVoters) Count() int {
	return len(tv.powers)
}

func (tv *topVoters) get(addr types.AccountID) *list.Element {
	return tv.powers[addr]
}

func (tv *topVoters) set(addr types.AccountID, e *list.Element) {
	tv.powers[addr] = e
}

func (tv *topVoters) getVotingPower(addr types.AccountID) *votingPower {
	if e := tv.powers[addr]; e != nil {
		return toVotingPower(e)
	}
	return nil
}

func (tv *topVoters) powerOf(addr types.AccountID) *big.Int {
	if vp := tv.getVotingPower(addr); vp != nil {
		return vp.getPower()
	}
	return nil
}

func (tv *topVoters) update(v *votingPower) (vp *votingPower) {
	var e *list.Element

	if e = tv.get(v.getAddr()); e != nil {
		vp = toVotingPower(e)
		vp.setPower(v.getPower())
		orderedListElemAdd(tv.list, e,
			func(e *list.Element) bool {
				existing := toVotingPower(e).getPower()
				curr := v.getPower()
				return tv.cmp(curr, existing)
			},
		)
	} else {
		vp = v
		e = orderedListAdd(tv.list, v,
			func(e *list.Element) bool {
				return tv.cmp(v.getPower(), toVotingPower(e).getPower())
			},
		)
		tv.set(v.getAddr(), e)
	}
	return
}

// Voters Power Ranking (VPR)
type vpr struct {
	changes    map[types.AccountID]*big.Int // temporary buffer for update
	voters     *topVoters
	store      *vprStore
	totalPower *big.Int
}

func (v *vpr) getTotalPower() *big.Int {
	return new(big.Int).Set(v.totalPower)
}

func (v *vpr) lowest() *votingPower {
	return v.store.getLowest()
}

func (v *vpr) updateLowest(vp *votingPower) {
	v.store.updateLowest(vp)
}

func newVpr() *vpr {
	return &vpr{
		voters:     newTopVoters(vprMax),
		store:      newVprStore(vprBucketsMax),
		totalPower: new(big.Int),
		changes:    make(map[types.AccountID]*big.Int, vprMax),
	}
}

func loadVpr(s dataGetter) (*vpr, error) {
	v := newVpr()

	for i := uint8(0); i < vprBucketsMax; i++ {
		var (
			vps []*votingPower
			err error
		)
		if vps, err = v.store.read(s, i); err != nil {
			return nil, err
		}
		for _, vp := range vps {
			v.store.addTail(i, vp)
			v.setVotingPower(vp)
		}
	}

	return v, nil
}

func (v *vpr) setVotingPower(vp *votingPower) *votingPower {
	rv := v.voters.update(vp)
	v.addTotal(rv.getPower())

	return rv
}

func (v *vpr) addTotal(delta *big.Int) {
	total := v.totalPower
	total.Add(total, delta)
}

func (v *vpr) votingPowerOf(addr types.AccountID) *big.Int {
	return v.voters.powerOf(addr)
}

func (v *vpr) updateRank(addr types.AccountID, delta *big.Int) *votingPower {
	vp := v.voters.getVotingPower(addr)
	if vp != nil {
		pwr := vp.getPower()
		pwr.Add(pwr, delta)
	} else {
		vp = newVotingPower(addr, delta)
	}

	rv := v.voters.update(vp)
	v.addTotal(delta)

	return rv
}

func (v *vpr) prepare(addr types.AccountID, fn func(lhs *big.Int)) {
	if _, exist := v.changes[addr]; !exist {
		v.changes[addr] = &big.Int{}
	}
	// Each entry of v.changes corresponds to the change (increment or
	// decrement) of voting power. It is added to later by calling the v.apply
	// method.
	ch := v.changes[addr]

	fn(ch)
}

func (v *vpr) add(addr types.AccountID, power *big.Int) {
	v.prepare(addr,
		func(lhs *big.Int) {
			lhs.Add(lhs, power)
		},
	)
}

func (v *vpr) sub(addr types.AccountID, power *big.Int) {
	v.prepare(addr,
		func(lhs *big.Int) {
			lhs.Sub(lhs, power)
		},
	)
}

func (v *vpr) apply(s *state.ContractState) (int, error) {
	var (
		nApplied = 0
		updRows  = make(map[uint8]interface{})
	)

	for addr, pow := range v.changes {
		if pow.Cmp(zeroValue) != 0 {
			vp := v.updateRank(addr, pow)

			if s != nil {
				i := v.store.update(vp)
				if _, exist := updRows[i]; !exist {
					updRows[i] = struct{}{}
				}
			}

			delete(v.changes, addr)
			// TODO: Remove a victim chosen above from the VPR bucket.
			nApplied++
		}
	}

	for i, _ := range updRows {
		if err := v.store.write(s, i); err != nil {
			return 0, err
		}
	}

	return nApplied, nil
}

func (v *vpr) Bingo(seed []byte) {
}

func vprKey(i uint8) []byte {
	var vk []byte = vprKeyPrefix
	return append(vk, []byte(fmt.Sprintf("%v", i))...)
}

func toVotingPower(e *list.Element) *votingPower {
	return e.Value.(*votingPower)
}
