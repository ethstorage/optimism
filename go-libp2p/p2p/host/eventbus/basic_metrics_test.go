package eventbus

import (
	"reflect"
	"testing"

	"github.com/libp2p/go-libp2p/core/event"
)

func BenchmarkEventEmitted(b *testing.B) {
	b.ReportAllocs()
	types := []reflect.Type{
		reflect.TypeOf(new(event.EvtLocalAddressesUpdated)),
		reflect.TypeOf(new(event.EvtNATDeviceTypeChanged)),
		reflect.TypeOf(new(event.EvtLocalProtocolsUpdated)),
	}
	mt := NewMetricsTracer()
	for i := 0; i < b.N; i++ {
		mt.EventEmitted(types[i%len(types)])
	}
}

func BenchmarkSubscriberQueueLength(b *testing.B) {
	b.ReportAllocs()
	names := []string{"s1", "s2", "s3", "s4"}
	mt := NewMetricsTracer()
	for i := 0; i < b.N; i++ {
		mt.SubscriberQueueLength(names[i%len(names)], i)
	}
}
