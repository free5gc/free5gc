package context_test

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/pkg/factory"
)

func TestUeIPPool(t *testing.T) {
	ueIPPool := context.NewUEIPPool(&factory.UEIPPool{
		Cidr: "10.10.0.0/24",
	})

	require.NotNil(t, ueIPPool)

	var allocIP net.IP

	// make allowed ip pools
	var ipPoolList []net.IP
	for i := 0; i <= 255; i += 1 {
		ipStr := fmt.Sprintf("10.10.0.%d", i)
		ipPoolList = append(ipPoolList, net.ParseIP(ipStr).To4())
	}

	// allocate
	for i := 0; i < 256; i += 1 {
		allocIP = ueIPPool.Allocate(nil)
		require.Contains(t, ipPoolList, allocIP)
	}

	// ip pool is empty
	allocIP = ueIPPool.Allocate(nil)
	require.Nil(t, allocIP)

	// release IP
	for _, i := range rand.Perm(256) {
		ueIPPool.Release(ipPoolList[i])
	}

	// allocate specify ip
	for _, ip := range ipPoolList {
		allocIP = ueIPPool.Allocate(ip)
		require.Equal(t, ip, allocIP)
	}
}

func TestUeIPPool_ExcludeRange(t *testing.T) {
	ueIPPool := context.NewUEIPPool(&factory.UEIPPool{
		Cidr: "10.10.0.0/24",
	})

	require.Equal(t, 0x0a0a0000, ueIPPool.Pool().Min())
	require.Equal(t, 0x0a0a00FF, ueIPPool.Pool().Max())
	require.Equal(t, 256, ueIPPool.Pool().Remain())

	excludeUeIPPool := context.NewUEIPPool(&factory.UEIPPool{
		Cidr: "10.10.0.0/28",
	})

	require.Equal(t, 0x0a0a0000, excludeUeIPPool.Pool().Min())
	require.Equal(t, 0x0a0a000F, excludeUeIPPool.Pool().Max())

	require.Equal(t, 16, excludeUeIPPool.Pool().Remain())

	err := ueIPPool.Exclude(excludeUeIPPool)
	require.NoError(t, err)
	require.Equal(t, 240, ueIPPool.Pool().Remain())

	for i := 16; i <= 255; i++ {
		allocate := ueIPPool.Allocate(nil)
		require.Equal(t, net.ParseIP(fmt.Sprintf("10.10.0.%d", i)).To4(), allocate)

		ueIPPool.Release(allocate)
	}
}
