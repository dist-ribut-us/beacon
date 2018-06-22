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
	b.router.Register(b)
	return b
}

// ServiceID for Beacon service
func (*Beacon) ServiceID() uint32 {
	return message.BeaconService
}

// NetQueryHandler for Beacon service
func (b *Beacon) NetQueryHandler(q ipcrouter.NetQuery) {
	switch t := q.GetType(); t {
	case message.GetIP:
		log.Info(log.Lbl("responding_to_get_ip"))
		q.Respond(q.GetAddrpb())
	default:
		log.Info(log.Lbl("unknown_net_message"), t)
	}
}

// requestOverlayPort then register Beacon with overlay as service
func (b *Beacon) requestOverlayPort() {
	b.router.RequestServicePort("Overlay", b.pool, func(r ipcrouter.Response) {
		b.overlay = rnet.Port(serial.UnmarshalUint16(r.GetBody()))
		log.Info(log.Lbl("overlay_port"), b.overlay)
		b.router.RegisterWithOverlay(message.BeaconService, b.overlay)
	})
}

// Run the beacon service
func (b *Beacon) Run() {
	go b.requestOverlayPort()
	b.router.Run()
}
