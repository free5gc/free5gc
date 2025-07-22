package context

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/free5gc/smf/internal/context/pool"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
)

// UeIPPool represent IPv4 address pool for UE
type UeIPPool struct {
	ueSubNet *net.IPNet
	pool     *pool.LazyReusePool
}

func NewUEIPPool(factoryPool *factory.UEIPPool) *UeIPPool {
	_, ipNet, err := net.ParseCIDR(factoryPool.Cidr)
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	minAddr, maxAddr, err := calcAddrRange(ipNet)
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	newPool, err := pool.NewLazyReusePool(int(minAddr), int(maxAddr))
	if err != nil {
		logger.InitLog.Errorln(err)
		return nil
	}

	ueIPPool := &UeIPPool{
		ueSubNet: ipNet,
		pool:     newPool,
	}
	return ueIPPool
}

func (ueIPPool *UeIPPool) Allocate(request net.IP) net.IP {
	var allocVal int
	var ok bool
	if request != nil {
		allocVal = int(binary.BigEndian.Uint32(request))
		ok = ueIPPool.pool.Use(allocVal)
		if !ok {
			logger.CtxLog.Warnf("IP[%s] is used in Pool[%+v]", request, ueIPPool.ueSubNet)
			return nil
		}
		// if allocated request IP address
		goto RETURNIP
	}

	allocVal, ok = ueIPPool.pool.Allocate()
	if !ok {
		logger.CtxLog.Warnf("Pool is empty: %+v", ueIPPool.ueSubNet)
		return nil
	}

RETURNIP:
	retIP := uint32ToIP(uint32(allocVal))
	logger.CtxLog.Infof("Allocated UE IP address: %s", retIP)
	return retIP
}

func (ueIPPool *UeIPPool) Exclude(excludePool *UeIPPool) error {
	excludeMin := excludePool.pool.Min()
	excludeMax := excludePool.pool.Max()
	if err := ueIPPool.pool.Reserve(excludeMin, excludeMax); err != nil {
		return fmt.Errorf("exclude uePool fail: %v", err)
	}
	return nil
}

func (u *UeIPPool) Pool() *pool.LazyReusePool {
	return u.pool
}

func uint32ToIP(intval uint32) net.IP {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, intval)
	return buf
}

func (ueIPPool *UeIPPool) Release(addr net.IP) {
	addrVal := binary.BigEndian.Uint32(addr)
	res := ueIPPool.pool.Free(int(addrVal))
	if !res {
		logger.CtxLog.Warnf("failed to release UE Address: %s", addr)
	}
	logger.CtxLog.Debug(ueIPPool.dump())
}

func (ueIPPool *UeIPPool) dump() string {
	str := "["
	elements := ueIPPool.pool.Dump()
	for index, element := range elements {
		var firstAddr net.IP
		var lastAddr net.IP
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[0]))
		firstAddr = buf
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(element[1]))
		lastAddr = buf
		if index > 0 {
			str += ("->")
		}
		str += fmt.Sprintf("{%s - %s}", firstAddr.String(), lastAddr.String())
	}
	str += ("]")
	return str
}

func isOverlap(pools []*UeIPPool) bool {
	if len(pools) < 2 {
		// no need to check
		return false
	}
	for i := 0; i < len(pools)-1; i++ {
		for j := i + 1; j < len(pools); j++ {
			if pools[i].pool.IsJoint(pools[j].pool) {
				return true
			}
		}
	}
	return false
}

func calcAddrRange(ipNet *net.IPNet) (minAddr, maxAddr uint32, err error) {
	maskVal := binary.BigEndian.Uint32(ipNet.Mask)
	baseIPVal := binary.BigEndian.Uint32(ipNet.IP)
	// move removing network and broadcast address later
	minAddr = (baseIPVal & maskVal)
	maxAddr = (baseIPVal | ^maskVal)
	if minAddr > maxAddr {
		return minAddr, maxAddr, errors.New("mask is invalid")
	}
	return minAddr, maxAddr, nil
}
