package pfcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/logger"
	logger_util "github.com/free5gc/util/logger"
)

func TestRemoteNode(t *testing.T) {
	t.Run("sess is not found before create", func(t *testing.T) {
		n := NewRemoteNode(
			"smf1",
			nil,
			&LocalNode{},
			forwarder.Empty{},
			logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf1"),
		)
		for i := 0; i < 3; i++ {
			_, err := n.Sess(uint64(i))
			assert.NotNil(t, err)
		}
	})

	t.Run("new multiple session", func(t *testing.T) {
		n := NewRemoteNode(
			"smf1",
			nil,
			&LocalNode{},
			forwarder.Empty{},
			logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf1"),
		)

		testcases := []struct {
			localID  uint64
			remoteID uint64
		}{
			{1, 10}, {2, 20}, {3, 30},
		}

		for _, tc := range testcases {
			sess := n.NewSess(tc.remoteID)
			assert.Equal(t, tc.localID, sess.LocalID)
			assert.Equal(t, tc.remoteID, sess.RemoteID)
		}

		// assure the session stored in the node
		for _, tc := range testcases {
			sess, err := n.Sess(tc.localID)
			assert.Nil(t, err)
			assert.Equal(t, tc.localID, sess.LocalID)
			assert.Equal(t, tc.remoteID, sess.RemoteID)
		}
	})

	t.Run("delete 0 no effect before create", func(t *testing.T) {
		n := NewRemoteNode(
			"smf1",
			nil,
			&LocalNode{},
			forwarder.Empty{},
			logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf1"),
		)
		report := n.DeleteSess(0)
		assert.Nil(t, report)
	})
	t.Run("delete should success after create", func(t *testing.T) {
		n := NewRemoteNode(
			"smf1",
			nil,
			&LocalNode{},
			forwarder.Empty{},
			logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf1"),
		)

		testcases := []struct {
			localID  uint64
			remoteID uint64
		}{
			{1, 10}, {2, 20}, {3, 30},
		}

		for _, tc := range testcases {
			n.NewSess(tc.remoteID)
		}

		for _, tc := range testcases {
			n.DeleteSess(tc.localID)
		}

		// assure the session is deleted
		for _, tc := range testcases {
			_, err := n.Sess(tc.localID)
			assert.NotNil(t, err)
		}

		// delete again should have no effect
		for _, tc := range testcases {
			report := n.DeleteSess(tc.localID)
			assert.Nil(t, report)
		}
	})
}

func TestRemoteNode_multipleSMF(t *testing.T) {
	var lnode LocalNode
	n1 := NewRemoteNode(
		"smf1",
		nil,
		&lnode,
		forwarder.Empty{},
		logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf1"),
	)
	n2 := NewRemoteNode(
		"smf2",
		nil,
		&lnode,
		forwarder.Empty{},
		logger.PfcpLog.WithField(logger_util.FieldControlPlaneNodeID, "smf2"),
	)
	t.Run("new smf1 r-SEID=10", func(t *testing.T) {
		sess := n1.NewSess(10)
		if sess.LocalID != 1 {
			t.Errorf("want 1; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 10 {
			t.Errorf("want 10; but got %v\n", sess.RemoteID)
		}
	})
	t.Run("new smf2 r-SEID=10", func(t *testing.T) {
		sess := n2.NewSess(10)
		if sess.LocalID != 2 {
			t.Errorf("want 2; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 10 {
			t.Errorf("want 10; but got %v\n", sess.RemoteID)
		}
	})
	t.Run("get smf1 l-SEID=1", func(t *testing.T) {
		sess, err := n1.Sess(1)
		if err != nil {
			t.Fatal(err)
		}
		if sess.LocalID != 1 {
			t.Errorf("want 1; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 10 {
			t.Errorf("want 10; but got %v\n", sess.RemoteID)
		}
	})
	t.Run("get smf2 l-SEID=2", func(t *testing.T) {
		sess, err := n2.Sess(2)
		if err != nil {
			t.Fatal(err)
		}
		if sess.LocalID != 2 {
			t.Errorf("want 2; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 10 {
			t.Errorf("want 10; but got %v\n", sess.RemoteID)
		}
	})
	t.Run("get smf1 l-SEID=2", func(t *testing.T) {
		_, err := n1.Sess(2)
		if err == nil {
			t.Errorf("want error; but not error")
		}
	})
	t.Run("get smf2 l-SEID=1", func(t *testing.T) {
		_, err := n2.Sess(1)
		if err == nil {
			t.Errorf("want error; but not error")
		}
	})
	t.Run("new smf1:20", func(t *testing.T) {
		sess := n1.NewSess(20)
		if sess.LocalID != 3 {
			t.Errorf("want 3; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 20 {
			t.Errorf("want 20; but got %v\n", sess.RemoteID)
		}
	})
	t.Run("get smf2 l-SEID=3", func(t *testing.T) {
		_, err := n2.Sess(3)
		if err == nil {
			t.Errorf("want error; but not error")
		}
	})
	t.Run("reset smf1", func(t *testing.T) {
		n1.Reset()
	})
	t.Run("get smf1 l-SEID=1", func(t *testing.T) {
		_, err := n1.Sess(1)
		if err == nil {
			t.Errorf("want error; but not error")
		}
	})
	t.Run("get smf1 l-SEID=3", func(t *testing.T) {
		_, err := n1.Sess(3)
		if err == nil {
			t.Errorf("want error; but not error")
		}
	})
	t.Run("get smf2 l-SEID=2", func(t *testing.T) {
		sess, err := n2.Sess(2)
		if err != nil {
			t.Fatal(err)
		}
		if sess.LocalID != 2 {
			t.Errorf("want 2; but got %v\n", sess.LocalID)
		}
		if sess.RemoteID != 10 {
			t.Errorf("want 10; but got %v\n", sess.RemoteID)
		}
	})
}

func TestLocalNode(t *testing.T) {
	t.Run("new session", func(t *testing.T) {
		lnode := LocalNode{}
		sess := lnode.NewSess(10, BUFFQ_LEN)
		assert.Equal(t, uint64(1), sess.LocalID)
		assert.Equal(t, uint64(10), sess.RemoteID)
	})

	t.Run("recycle LocalID", func(t *testing.T) {
		lnode := LocalNode{
			sess: []*Sess{},
			free: []uint64{},
		}
		sess := lnode.NewSess(10, BUFFQ_LEN)
		recycleLocalID := 1
		assert.Equal(t, uint64(recycleLocalID), sess.LocalID)
		assert.Equal(t, uint64(10), sess.RemoteID)
	})
}
