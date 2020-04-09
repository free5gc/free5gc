package ike_handler

import (
	"errors"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_ike/ike_message"

	"github.com/vishvananda/netlink"
)

type XFRMEncryptionAlgorithmType uint16

func (xfrmEncryptionAlgorithmType XFRMEncryptionAlgorithmType) String() string {
	switch xfrmEncryptionAlgorithmType {
	case ike_message.ENCR_DES:
		return "cbc(des)"
	case ike_message.ENCR_3DES:
		return "cbc(des3_ede)"
	case ike_message.ENCR_CAST:
		return "cbc(cast5)"
	case ike_message.ENCR_BLOWFISH:
		return "cbc(blowfish)"
	case ike_message.ENCR_NULL:
		return "ecb(cipher_null)"
	case ike_message.ENCR_AES_CBC:
		return "cbc(aes)"
	case ike_message.ENCR_AES_CTR:
		return "rfc3686(ctr(aes))"
	default:
		return ""
	}
}

type XFRMIntegrityAlgorithmType uint16

func (xfrmIntegrityAlgorithmType XFRMIntegrityAlgorithmType) String() string {
	switch xfrmIntegrityAlgorithmType {
	case ike_message.AUTH_HMAC_MD5_96:
		return "hmac(md5)"
	case ike_message.AUTH_HMAC_SHA1_96:
		return "hmac(sha1)"
	case ike_message.AUTH_AES_XCBC_96:
		return "xcbc(aes)"
	default:
		return ""
	}
}

func ApplyXFRMRule(n3iwf_is_initiator bool, childSecurityAssociation *n3iwf_context.ChildSecurityAssociation) error {
	N3IWFSelf := n3iwf_context.N3IWFSelf()

	// Build XFRM information data structure for incoming traffic.

	// Mark
	mark := &netlink.XfrmMark{
		Value: N3IWFSelf.Mark,
	}

	// Direction: {private_network} -> this_server
	// State
	var xfrmEncryptionAlgorithm, xfrmIntegrityAlgorithm *netlink.XfrmStateAlgo
	if n3iwf_is_initiator {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: XFRMEncryptionAlgorithmType(childSecurityAssociation.EncryptionAlgorithm).String(),
			Key:  childSecurityAssociation.ResponderToInitiatorEncryptionKey,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
				Key:  childSecurityAssociation.ResponderToInitiatorIntegrityKey,
			}
		}
	} else {
		xfrmEncryptionAlgorithm = &netlink.XfrmStateAlgo{
			Name: XFRMEncryptionAlgorithmType(childSecurityAssociation.EncryptionAlgorithm).String(),
			Key:  childSecurityAssociation.InitiatorToResponderEncryptionKey,
		}
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm = &netlink.XfrmStateAlgo{
				Name: XFRMIntegrityAlgorithmType(childSecurityAssociation.IntegrityAlgorithm).String(),
				Key:  childSecurityAssociation.InitiatorToResponderIntegrityKey,
			}
		}
	}

	xfrmState := new(netlink.XfrmState)

	xfrmState.Src = childSecurityAssociation.PeerPublicIPAddr
	xfrmState.Dst = childSecurityAssociation.LocalPublicIPAddr
	xfrmState.Proto = netlink.XFRM_PROTO_ESP
	xfrmState.Mode = netlink.XFRM_MODE_TUNNEL
	xfrmState.Spi = int(childSecurityAssociation.SPI)
	xfrmState.Mark = mark
	xfrmState.Auth = xfrmIntegrityAlgorithm
	xfrmState.Crypt = xfrmEncryptionAlgorithm
	xfrmState.ESN = childSecurityAssociation.ESN

	// Commit xfrm state to netlink
	var err error
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		ikeLog.Errorf("[IKE] Set XFRM rules failed: %+v", err)
		return errors.New("Set XFRM state rule failed")
	}

	// Policy
	xfrmPolicyTemplate := netlink.XfrmPolicyTmpl{
		Src:   xfrmState.Src,
		Dst:   xfrmState.Dst,
		Proto: xfrmState.Proto,
		Mode:  xfrmState.Mode,
		Spi:   xfrmState.Spi,
	}

	xfrmPolicy := new(netlink.XfrmPolicy)

	xfrmPolicy.Src = &childSecurityAssociation.TrafficSelectorRemote
	xfrmPolicy.Dst = &childSecurityAssociation.TrafficSelectorLocal
	xfrmPolicy.Proto = netlink.Proto(childSecurityAssociation.SelectedIPProtocol)
	xfrmPolicy.Dir = netlink.XFRM_DIR_IN
	xfrmPolicy.Mark = mark
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		ikeLog.Errorf("[IKE] Set XFRM rules failed: %+v", err)
		return errors.New("Set XFRM policy rule failed")
	}

	// Direction: this_server -> {private_network}
	// State
	if n3iwf_is_initiator {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.InitiatorToResponderEncryptionKey
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.InitiatorToResponderIntegrityKey
		}
	} else {
		xfrmEncryptionAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorEncryptionKey
		if childSecurityAssociation.IntegrityAlgorithm != 0 {
			xfrmIntegrityAlgorithm.Key = childSecurityAssociation.ResponderToInitiatorIntegrityKey
		}
	}

	xfrmState.Src, xfrmState.Dst = xfrmState.Dst, xfrmState.Src

	// Commit xfrm state to netlink
	if err = netlink.XfrmStateAdd(xfrmState); err != nil {
		ikeLog.Errorf("[IKE] Set XFRM rules failed: %+v", err)
		return errors.New("Set XFRM state rule failed")
	}

	// Policy
	xfrmPolicyTemplate.Src, xfrmPolicyTemplate.Dst = xfrmPolicyTemplate.Dst, xfrmPolicyTemplate.Src

	xfrmPolicy.Src, xfrmPolicy.Dst = xfrmPolicy.Dst, xfrmPolicy.Src
	xfrmPolicy.Dir = netlink.XFRM_DIR_OUT
	xfrmPolicy.Tmpls = []netlink.XfrmPolicyTmpl{
		xfrmPolicyTemplate,
	}

	// Commit xfrm policy to netlink
	if err = netlink.XfrmPolicyAdd(xfrmPolicy); err != nil {
		ikeLog.Errorf("[IKE] Set XFRM rules failed: %+v", err)
		return errors.New("Set XFRM policy rule failed")
	}

	return nil
}
