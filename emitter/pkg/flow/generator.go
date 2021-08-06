package flow

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type RndGenerator struct {
	sequence     int32
	podAddresses []string
}

const flowTemplate = `{"Type":"IPFix","TimeReceived":"%d","SequenceNum":"%d","SamplingRate":"0","SamplerAddress":"0.0.0.0",` +
	`"TimeFlowStart":"%d","TimeFlowEnd":"%d","Bytes":"%d","Packets":"%d","SrcAddr":%q,"DstAddr":%q,` +
	`"Etype":"2048","Proto":"6","SrcPort":"%d","DstPort":"%d",` +
	`"InIf":"11","OutIf":"65533","SrcMac":"00:00:00:00:00:00","DstMac":"00:00:00:00:00:00",` +
	`"SrcVlan":"0","DstVlan":"0","VlanId":"0","IngressVrfID":"0","EgressVrfID":"0","IPTos":"0",` +
	`"ForwardingStatus":"0","IPTTL":"0","TCPFlags":"16","IcmpType":"0","IcmpCode":"0",` +
	`"IPv6FlowLabel":"0","FragmentId":"0","FragmentOffset":"0","BiFlowDirection":"0","SrcAS":"0",` +
	`"DstAS":"0","NextHop":"0.0.0.0","NextHopAS":"0","SrcNet":"0","DstNet":"0","HasEncap":"false",` +
	`"SrcAddrEncap":"<nil>","DstAddrEncap":"<nil>","ProtoEncap":"0","EtypeEncap":"0",` +
	`"IPTosEncap":"0","IPTTLEncap":"0","IPv6FlowLabelEncap":"0","FragmentIdEncap":"0",` +
	`"FragmentOffsetEncap":"0","HasMPLS":"false","MPLSCount":"0","MPLS1TTL":"0",` +
	`"MPLS1Label":"0","MPLS2TTL":"0","MPLS2Label":"0","MPLS3TTL":"0","MPLS3Label":"0",` +
	`"MPLSLastTTL":"0","MPLSLastLabel":"0","HasPPP":"false","PPPAddressControl":"0",` +
	`"K8SSrcPodName":%q,"K8SSrcPodNamespace":"namespace","K8SSrcPodNode":"0.0.0.0",` +
	`"K8SDstPodName":%q,"K8SDstPodNamespace":"namespace","K8SDstPodNode":"0.0.0.0"}`

func NewRndGenerator(podsBaseAddr string, numPods int) RndGenerator {
	return RndGenerator{
		sequence:     0,
		podAddresses: ipRange(podsBaseAddr, numPods),
	}
}

func (g *RndGenerator) Generate() (payload, srcPod, dstPod string) {
	now := time.Now().Unix()
	g.sequence++
	srcNum := rand.Intn(len(g.podAddresses))
	dstNum := rand.Intn(len(g.podAddresses))
	srcPod = g.podAddresses[srcNum]
	dstPod = g.podAddresses[dstNum]
	payload = fmt.Sprintf(
		flowTemplate, now, g.sequence,
		now, now, rand.Intn(4095)+1, rand.Intn(63)+1, srcPod, dstPod,
		rand.Intn(65000)+1, rand.Intn(65000)+1,
		srcPod, dstPod)
	return payload, srcPod, dstPod
}

func ipRange(base string, length int) []string {
	rng := make([]string, 0, length)
	ip := binary.BigEndian.Uint32(net.ParseIP(base).To4())
	for i := 0; i < length; i++ {
		ip++
		newIP := make(net.IP, 4)
		binary.BigEndian.PutUint32(newIP, ip)
		rng = append(rng, newIP.String())
	}
	return rng
}
