package network

import (
	"github.com/CasperHollemans/PepperChain/internal/logging"
)

type Network struct {
	Nodes []string
}

type NetworkService interface {
	RegisterNode(address string)
	GetNodes() []string
}

type NetworkServiceImpl struct {
	logger  logging.Logger
	network Network
}

func NewNetworkService(logger logging.Logger) *NetworkServiceImpl {
	return &NetworkServiceImpl{
		network: Network{},
		logger:  logger,
	}
}

func (n *NetworkServiceImpl) RegisterNode(address string) {
	n.logger.Info("Registering node", address)
	n.network.Nodes = append(n.network.Nodes, address)
}

func (n *NetworkServiceImpl) GetNodes() []string {
	return n.network.Nodes
}
