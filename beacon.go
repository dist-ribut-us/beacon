package beacon

import (
	"github.com/dist-ribut-us/ipcrouter"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/message"
	"github.com/dist-ribut-us/rnet"
	"github.com/dist-ribut-us/serial"
)

// Beacon service struct
type Beacon struct {
	router  *ipcrouter.Router
	overlay rnet.Port
	pool    rnet.Port
}

// New Beacon struct
func New(router *ipcrouter.Router, pool rnet.Port) *Beacon {
	b := &Beacon{
		router: router,
		pool:   pool,
	}
	b.router.Register(message.BeaconService, b.handler)
	return b
}

// Run listens on the IPC channel and handles any messages it receives.
func (b *Beacon) handler(bs *ipcrouter.Base) {
	if bs.IsFromNet() {
		b.handleNetReceive(bs)
	} else {
		log.Info(log.Lbl("unknown_type"), bs.GetType())
	}
}

func (b *Beacon) handleNetReceive(bs *ipcrouter.Base) {
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
	b.router.RequestServicePort("Overlay", b.pool, func(r *ipcrouter.Base) {
		b.overlay = rnet.Port(serial.UnmarshalUint16(r.Body))
		log.Info(log.Lbl("overlay_port"), b.overlay)
		b.router.RegisterWithOverlay(message.BeaconService, b.overlay)
	})
}

// Run the beacon service
func (b *Beacon) Run() {
	go b.requestOverlayPort()
	b.router.Run()
}
