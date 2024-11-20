// Package env TODO
package env

import "net/url"

// const The default dmp backend service address for each environment
const (
	DefaultDMPAddrPrd  = "https://openapi.abetterchoice.ai"
	DefaultDMPAddrTest = "https://openapi.abetterchoice.ai"
)

// dmpAddrIndex DMP backend address index for each environment, replace the default address through RegisterDMPAddr
var dmpAddrIndex = map[Type]string{
	TypePrd:  DefaultDMPAddrPrd,
	TypeTest: DefaultDMPAddrTest,
}

// RegisterDMPAddr Register the dmp backend address in the specified environment
func RegisterDMPAddr(envType Type, addr string) error {
	_, err := url.Parse(addr)
	if err != nil {
		return err
	}
	dmpAddrIndex[envType] = addr
	return nil
}

// GetDMPAddr Get the dmp backend address in the specified environment.
// If the specified environment configuration does not exist, the official environment address is used by default.
func GetDMPAddr(envType Type) string {
	addr, ok := dmpAddrIndex[envType]
	if !ok {
		return dmpAddrIndex[TypePrd]
	}
	return addr
}
