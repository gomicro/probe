package cmd

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gomicro/bogus"
	. "github.com/onsi/gomega"
)

func TestProbe(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Probe", func() {
		g.It("should probe an endpoint", func() {
			server := bogus.New()
			server.AddPath("/v1/status").
				SetMethods("GET").
				SetPayload([]byte("ok"))
			h, p := server.HostPort()

			args := []string{fmt.Sprintf("http://%v:%v/v1/status", h, p)}
			probe(nil, args)

			hr := server.HitRecords()
			Expect(len(hr)).To(Equal(1))
			Expect(hr[0].Verb).To(Equal("GET"))
		})
	})
}
