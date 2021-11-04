package ssi

import (
	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	log "github.com/sirupsen/logrus"
)

// CosmosRegistry vdr registry.
type CosmosRegistry struct{}

func (r *CosmosRegistry) Resolve(did string, opts ...vdr.DIDMethodOption) (*did.DocResolution, error) {
	log.Infoln("VDR Resolve Called", did)
	return nil, nil
}
func (r *CosmosRegistry) Create(method string, did *did.Doc, opts ...vdr.DIDMethodOption) (*did.DocResolution, error) {
	log.Infoln("VDR Create Called", did.ID)
	return nil, nil
}
func (r *CosmosRegistry) Update(did *did.Doc, opts ...vdr.DIDMethodOption) error {
	log.Infoln("VDR Update Called", did)
	return nil
}
func (r *CosmosRegistry) Deactivate(did string, opts ...vdr.DIDMethodOption) error {
	log.Infoln("VDR Resolve Deactivate", did)
	return nil
}
func (r *CosmosRegistry) Close() error {
	log.Infoln("VDR Resolve Close")
	return nil
}

// CosmosVDR verifiable data registry
type CosmosVDR struct {
}

func (vdr CosmosVDR) Read(did string, opts ...vdr.DIDMethodOption) (*did.DocResolution, error) {
	log.Infoln("VDR Read Called", did)
	return nil, nil
}
func (vdr CosmosVDR) Create(did *did.Doc, opts ...vdr.DIDMethodOption) (*did.DocResolution, error) {
	log.Infoln("VDR Create Called", did.ID)
	return nil, nil
}
func (vdr CosmosVDR) Accept(method string) bool {
	log.Infoln("VDR Accept Called", method)
	return true
}
func (vdr CosmosVDR) Update(did *did.Doc, opts ...vdr.DIDMethodOption) error   { return nil }
func (vdr CosmosVDR) Deactivate(did string, opts ...vdr.DIDMethodOption) error { return nil }
func (vdr CosmosVDR) Close() error                                             { return nil }
