package beacon

import (
	"github.com/dist-ribut-us/ipc"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/message"
	"github.com/dist-ribut-us/rnet"
	"github.com/dist-ribut-us/serial"
)

// BeaconID is the service ID for DHT
const BeaconID uint32 = 3819762595

// Beacon service struct
type Beacon struct {
	ipc     *ipc.Proc
	overlay rnet.Port
	pool    rnet.Port
}

// New Beacon struct
func New(proc *ipc.Proc, pool rnet.Port) *Beacon {
	b := &Beacon{
		ipc:  proc,
		pool: pool,
	}
	b.ipc.Handler(b.handler)
	return b
}

// Run listens on the IPC channel and handles any messages it receives.
func (b *Beacon) handler(bs *ipc.Base) {
	if bs.IsFromNet() {
		b.handleNetReceive(bs)
	} else {
		log.Info(log.Lbl("unknown_type"), bs.GetType())
	}
}

func (b *Beacon) handleNetReceive(bs *ipc.Base) {
	switch t := bs.GetType(); t {
	case message.GetIP:
		if bs.IsQuery() {
			log.Info(log.Lbl("responding_to_get_ip"))
			bs.Respond(bs.Addrpb)
		} else {
			log.Info(log.Lbl("non_query_get_ip"))
		}
	default:
		log.Info(log.Lbl("unknown_net_message"), t)
	}
}

// requestOverlayPort then register Beacon with overlay as service
func (b *Beacon) requestOverlayPort() {
	b.ipc.RequestServicePort("Overlay", b.pool, func(r *ipc.Base) {
		b.overlay = rnet.Port(serial.UnmarshalUint16(r.Body))
		log.Info(log.Lbl("overlay_port"), b.overlay)
		b.ipc.RegisterWithOverlay(BeaconID, b.overlay, nil)
	})
}

// Run the beacon service
func (b *Beacon) Run() {
	go b.requestOverlayPort()
	b.ipc.Run()
}
