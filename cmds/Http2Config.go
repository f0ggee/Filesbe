package cmds

import "net/http"

type Http2Config struct {
	NewConfig http.HTTP2Config
}

func ConfigHttp2() *Http2Config {

	NewConfig := http.HTTP2Config{
		MaxConcurrentStreams:          60,
		MaxDecoderHeaderTableSize:     0,
		MaxEncoderHeaderTableSize:     0,
		MaxReadFrameSize:              0,
		MaxReceiveBufferPerConnection: 0,
		MaxReceiveBufferPerStream:     0,
		SendPingTimeout:               0,
		PingTimeout:                   0,
		WriteByteTimeout:              0,
		PermitProhibitedCipherSuites:  false,
		CountError:                    nil,
	}
	return &Http2Config{NewConfig}
}
