package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/franela/goblin"
	"github.com/gomicro/bogus"
	. "github.com/onsi/gomega"
)

func TestProbe(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Probe", func() {
		var server *bogus.Bogus
		var host string

		g.BeforeEach(func() {
			server = bogus.New()
			h, p := server.HostPort()
			host = fmt.Sprintf("http://%v:%v/v1/status", h, p)
		})

		g.It("should probe an endpoint", func() {
			server.AddPath("/v1/status").
				SetMethods("GET").
				SetPayload([]byte("ok"))

			err := probeHttp(host)
			Expect(err).NotTo(HaveOccurred())

			hr := server.HitRecords()
			Expect(len(hr)).To(Equal(1))
			Expect(hr[0].Verb).To(Equal("GET"))
		})

		g.It("should error when a bad status is returned", func() {
			server.AddPath("/v1/status").
				SetMethods("GET").
				SetStatus(http.StatusNotFound).
				SetPayload([]byte("bad"))

			err := probeHttp(host)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrHttpBadStatus)).To(BeTrue())

			hr := server.HitRecords()
			Expect(len(hr)).To(Equal(1))
			Expect(hr[0].Verb).To(Equal("GET"))
		})
	})
}
