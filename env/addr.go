// Package env TODO
package env

import "net/url"

// const The default backend address for each environment
const (
	DefaultAddrPrd  = "https://mobile.abetterchoice.ai"
	DefaultAddrTest = "https://mobile.abetterchoice.ai"
)

// addrIndex
var addrIndex = map[Type]string{
	TypePrd:  DefaultAddrPrd,
	TypeTest: DefaultAddrTest,
}

// RegisterAddr Register backend address
func RegisterAddr(envType Type, addr string) error {
	_, err := url.Parse(addr)
	if err != nil {
		return err
	}
	addrIndex[envType] = addr
	return nil
}

// GetAddr Get the backend address in the specified environment.
// If the specified environment configuration does not exist, the official environment address is used by default.
func GetAddr(envType Type) string {
	addr, ok := addrIndex[envType]
	if !ok {
		return addrIndex[TypePrd]
	}
	return addr
}

var (
	// CacheServerSocket5Addr Cache service socket5 proxy address, set before SDK Init to take effect
	CacheServerSocket5Addr = ""
	// DMPServerSocket5Addr DMP user profile socket5 proxy address, which will take effect only if set before SDK Init
	DMPServerSocket5Addr = ""
)
