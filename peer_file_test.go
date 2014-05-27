package thermal

import (
	"testing"
)

// variation of seeds.json from telehash-c project
var seedsList1 = `
{
    "dca549c98b94197e79dfa2a7e4ad2e74fb144cda3f4710d4f40e2c75d975272e": {
        "paths": [
            {
                "type": "http",
                "http": "http://192.168.0.36:42424"
            },
            {
                "type": "ipv4",
                "ip": "127.0.0.1",
                "port": 42424
            },
            {
                "type": "ipv6",
                "ip": "fe80::bae8:56ff:fe43:3de4",
                "port": 42424
            }
        ],
        "parts": {
            "3a": "f0d2bfc8590a7e0016ce85dbf0f8f1883fb4f3dcc4701eab12ef83f972a2b87f",
            "2a": "0cb4f6137a745f1af2d31707550c03b99083180f6e69ec37918c220ecfa2972f",
            "1a": "b5a96d25802b3600ea99774138a650d5d1fa1f3cf3cb10ae8f1c58a527d85086"
        },
        "keys": {
            "3a": "MC5dfSfrAVCSugX75JbgVWtvCbxPqwLDUkc9TcS/qxE=",
            "2a": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqr12tXnpn707llkZfEcspB/D6KTcZM765+SnI5Z8JWkjc0Mrz9qZBB2YFLr2NmgCx0oLfSetmuHBNTT54sIAxQ/vxyykcMNGsSFg4WKhbsQXSrX4qChbhpIqMJkKa4mYZIb6qONA76G5/431u4+1sBRvfY0ewHChqGh0oThcaa50nT68f8ohIs1iUFm+SL8L9UL/oKN3Yg6drBYwpJi2Ex5Idyu4YQJwZ9sAQU49Pfs+LqhkHOascTmaa3+kTyTnp2iJ9wEuPg+AR3PJwxXnwYoWbH+Wr8gY6iLe0FQe8jXk6eLw9mqOhUcah8338MC83zSQcZriGVMq8qaQz0L9nwIDAQAB",
            "1a": "z6yCAC7r5XIr6C4xdxeX7RlSmGu9Xe73L1gv8qecm4/UEZAKR5iCxA=="
        }
    },
    "11111111dca549c98b94197e79dfa2a7e4ad2e74fb144cda3f4710d4f40e2c75d975272e": {
        "paths": [
            {
                "type": "ipv4",
                "ip": "127.0.0.1",
                "port": 42424
            }
        ],
        "parts": {
            "3a": "f0d2bfc8590a7e0016ce85dbf0f8f1883fb4f3dcc4701eab12ef83f972a2b87f"
        },
        "keys": {
            "3a": "MC5dfSfrAVCSugX75JbgVWtvCbxPqwLDUkc9TcS/qxE="
        }
    }
}
`

func TestLoadPeersFromJson(t *testing.T) {

	peers := loadPeersFromJson([]byte(seedsList1))
	for _, peer := range peers {
		t.Logf("peer: %s", peer)
	}

}
