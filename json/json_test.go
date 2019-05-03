package json

import (
	"testing"
)

func TestUnmarshalToContext(t *testing.T) {
	raw := `
{
  "ignition": { "version": "2.2.0" },
  "systemd": {
    "units": [{
      "name": "systemd-networkd.service",
      "dropins": [{
        "name": "debug.conf",
        "contents": "[Service]\nEnvironment=SYSTEMD_LOG_LEVEL=debug"
      }]
    }]
  }
}`
	res, err := UnmarshalToContext([]byte(raw))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("%+v", res)
}
