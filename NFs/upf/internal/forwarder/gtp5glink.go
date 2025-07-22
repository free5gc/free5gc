package forwarder

import (
	"net"
	"os"
	"syscall"

	"github.com/khirono/go-nl"
	"github.com/khirono/go-rtnllink"
	"github.com/khirono/go-rtnlroute"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/go-gtp5gnl"
)

type Gtp5gLink struct {
	mux    *nl.Mux
	rtconn *nl.Conn
	client *nl.Client
	link   *gtp5gnl.Link
	conn   *net.UDPConn
	f      *os.File
	log    *logrus.Entry
}

func OpenGtp5gLink(mux *nl.Mux, addr string, mtu uint32, log *logrus.Entry) (*Gtp5gLink, error) {
	g := &Gtp5gLink{
		log: log,
	}

	g.mux = mux

	rtconn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		return nil, errors.Wrap(err, "open")
	}
	g.rtconn = rtconn
	g.client = nl.NewClient(rtconn, mux)

	laddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "resolve addr")
	}
	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "listen")
	}
	g.conn = conn

	// TODO: Duplicate fd
	f, err := conn.File()
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "file")
	}
	g.f = f

	linkinfo := &nl.Attr{
		Type: syscall.IFLA_LINKINFO,
		Value: nl.AttrList{
			{
				Type:  rtnllink.IFLA_INFO_KIND,
				Value: nl.AttrString("gtp5g"),
			},
			{
				Type: rtnllink.IFLA_INFO_DATA,
				Value: nl.AttrList{
					{
						Type:  gtp5gnl.IFLA_FD1,
						Value: nl.AttrU32(f.Fd()),
					},
					{
						Type:  gtp5gnl.IFLA_HASHSIZE,
						Value: nl.AttrU32(131072),
					},
				},
			},
		},
	}
	attrs := []*nl.Attr{linkinfo}

	if mtu != 0 {
		attrs = append(attrs, &nl.Attr{
			Type:  syscall.IFLA_MTU,
			Value: nl.AttrU32(mtu),
		})
	}

	err = rtnllink.Create(g.client, "upfgtp", attrs...)
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "create")
	}
	err = rtnllink.Up(g.client, "upfgtp")
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "up")
	}
	link, err := gtp5gnl.GetLink("upfgtp")
	if err != nil {
		g.Close()
		return nil, errors.Wrap(err, "get link")
	}
	g.link = link
	return g, nil
}

func (g *Gtp5gLink) Close() {
	if g.f != nil {
		err := g.f.Close()
		if err != nil {
			g.log.Warnf("file close err: %+v", err)
		}
	}
	if g.conn != nil {
		err := g.conn.Close()
		if err != nil {
			g.log.Warnf("conn close err: %+v", err)
		}
	}
	if g.link != nil {
		err := rtnllink.Remove(g.client, "upfgtp")
		if err != nil {
			g.log.Warnf("rtnllink remove err: %+v", err)
		}
	}
	if g.rtconn != nil {
		g.rtconn.Close()
	}
}

func (g *Gtp5gLink) RouteAdd(dst *net.IPNet) error {
	r := &rtnlroute.Request{
		Header: rtnlroute.Header{
			Table:    syscall.RT_TABLE_MAIN,
			Scope:    syscall.RT_SCOPE_UNIVERSE,
			Protocol: syscall.RTPROT_STATIC,
			Type:     syscall.RTN_UNICAST,
		},
	}
	err := r.AddDst(dst)
	if err != nil {
		return err
	}
	err = r.AddIfName(g.link.Name)
	if err != nil {
		return err
	}
	return rtnlroute.Create(g.client, r)
}

func (g *Gtp5gLink) WriteTo(b []byte, addr net.Addr) (int, error) {
	return g.conn.WriteTo(b, addr)
}
