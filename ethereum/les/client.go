// Copyright 2016 The mouse Authors
// This file is part of the mouse library.
//
// The mouse library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The mouse library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the mouse library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Mouse Subprotocol.
package les

import (
	"fmt"
	"time"

	"github.com/marcopoloprotoco/mouse/accounts"
	"github.com/marcopoloprotoco/mouse/common"
	"github.com/marcopoloprotoco/mouse/common/hexutil"
	"github.com/marcopoloprotoco/mouse/common/mclock"
	"github.com/marcopoloprotoco/mouse/consensus"
	"github.com/marcopoloprotoco/mouse/core"
	"github.com/marcopoloprotoco/mouse/core/bloombits"
	"github.com/marcopoloprotoco/mouse/core/rawdb"
	"github.com/marcopoloprotoco/mouse/core/types"
	"github.com/marcopoloprotoco/mouse/mos"
	"github.com/marcopoloprotoco/mouse/mos/downloader"
	"github.com/marcopoloprotoco/mouse/mos/filters"
	"github.com/marcopoloprotoco/mouse/mos/gasprice"
	"github.com/marcopoloprotoco/mouse/event"
	"github.com/marcopoloprotoco/mouse/internal/ethapi"
	lpc "github.com/marcopoloprotoco/mouse/les/lespay/client"
	"github.com/marcopoloprotoco/mouse/light"
	"github.com/marcopoloprotoco/mouse/log"
	"github.com/marcopoloprotoco/mouse/node"
	"github.com/marcopoloprotoco/mouse/p2p"
	"github.com/marcopoloprotoco/mouse/p2p/enode"
	"github.com/marcopoloprotoco/mouse/params"
	"github.com/marcopoloprotoco/mouse/rpc"
)

type LightEthereum struct {
	lesCommons

	peers          *serverPeerSet
	reqDist        *requestDistributor
	retriever      *retrieveManager
	odr            *LesOdr
	relay          *lesTxRelay
	handler        *clientHandler
	txPool         *light.TxPool
	blockchain     *light.LightChain
	serverPool     *serverPool
	valueTracker   *lpc.ValueTracker
	dialCandidates enode.Iterator
	pruner         *pruner

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	ApiBackend     *LesApiBackend
	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager
	netRPCService  *ethapi.PublicNetAPI

	p2pServer *p2p.Server
}

// New creates an instance of the light client.
func New(stack *node.Node, config *mos.Config) (*LightEthereum, error) {
	chainDb, err := stack.OpenDatabase("lightchaindata", config.DatabaseCache, config.DatabaseHandles, "mos/db/chaindata/")
	if err != nil {
		return nil, err
	}
	lespayDb, err := stack.OpenDatabase("lespay", 0, 0, "mos/db/lespay")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newServerPeerSet()
	leth := &LightEthereum{
		lesCommons: lesCommons{
			genesis:     genesisHash,
			config:      config,
			chainConfig: chainConfig,
			iConfig:     light.DefaultClientIndexerConfig,
			chainDb:     chainDb,
			closeCh:     make(chan struct{}),
		},
		peers:          peers,
		eventMux:       stack.EventMux(),
		reqDist:        newRequestDistributor(peers, &mclock.System{}),
		accountManager: stack.AccountManager(),
		engine:         mos.CreateConsensusEngine(stack, chainConfig, &config.Ethash, nil, false, chainDb),
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   mos.NewBloomIndexer(chainDb, params.BloomBitsBlocksClient, params.HelperTrieConfirmations),
		valueTracker:   lpc.NewValueTracker(lespayDb, &mclock.System{}, requestList, time.Minute, 1/float64(time.Hour), 1/float64(time.Hour*100), 1/float64(time.Hour*1000)),
		p2pServer:      stack.Server(),
	}
	peers.subscribe((*vtSubscription)(leth.valueTracker))

	dnsdisc, err := leth.setupDiscovery(&stack.Config().P2P)
	if err != nil {
		return nil, err
	}
	leth.serverPool = newServerPool(lespayDb, []byte("serverpool:"), leth.valueTracker, dnsdisc, time.Second, nil, &mclock.System{}, config.UltraLightServers)
	peers.subscribe(leth.serverPool)
	leth.dialCandidates = leth.serverPool.dialIterator

	leth.retriever = newRetrieveManager(peers, leth.reqDist, leth.serverPool.getTimeout)
	leth.relay = newLesTxRelay(peers, leth.retriever)

	leth.odr = NewLesOdr(chainDb, light.DefaultClientIndexerConfig, leth.retriever)
	leth.chtIndexer = light.NewChtIndexer(chainDb, leth.odr, params.CHTFrequency, params.HelperTrieConfirmations, config.LightNoPrune)
	leth.bloomTrieIndexer = light.NewBloomTrieIndexer(chainDb, leth.odr, params.BloomBitsBlocksClient, params.BloomTrieFrequency, config.LightNoPrune)
	leth.odr.SetIndexers(leth.chtIndexer, leth.bloomTrieIndexer, leth.bloomIndexer)

	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}
	// Note: NewLightChain adds the trusted checkpoint so it needs an ODR with
	// indexers already set but not started yet
	if leth.blockchain, err = light.NewLightChain(leth.odr, leth.chainConfig, leth.engine, checkpoint); err != nil {
		return nil, err
	}
	leth.chainReader = leth.blockchain
	leth.txPool = light.NewTxPool(leth.chainConfig, leth.blockchain, leth.relay)

	// Set up checkpoint oracle.
	leth.oracle = leth.setupOracle(stack, genesisHash, config)

	// Note: AddChildIndexer starts the update process for the child
	leth.bloomIndexer.AddChildIndexer(leth.bloomTrieIndexer)
	leth.chtIndexer.Start(leth.blockchain)
	leth.bloomIndexer.Start(leth.blockchain)

	// Start a light chain pruner to delete useless historical data.
	leth.pruner = newPruner(chainDb, leth.chtIndexer, leth.bloomTrieIndexer)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		leth.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	leth.ApiBackend = &LesApiBackend{stack.Config().ExtRPCEnabled(), leth, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	leth.ApiBackend.gpo = gasprice.NewOracle(leth.ApiBackend, gpoParams)

	leth.handler = newClientHandler(config.UltraLightServers, config.UltraLightFraction, checkpoint, leth)
	if leth.handler.ulc != nil {
		log.Warn("Ultra light client is enabled", "trustedNodes", len(leth.handler.ulc.keys), "minTrustedFraction", leth.handler.ulc.fraction)
		leth.blockchain.DisableCheckFreq()
	}

	leth.netRPCService = ethapi.NewPublicNetAPI(leth.p2pServer, leth.config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(leth.APIs())
	stack.RegisterProtocols(leth.Protocols())
	stack.RegisterLifecycle(leth)

	return leth, nil
}

// vtSubscription implements serverPeerSubscriber
type vtSubscription lpc.ValueTracker

// registerPeer implements serverPeerSubscriber
func (v *vtSubscription) registerPeer(p *serverPeer) {
	vt := (*lpc.ValueTracker)(v)
	p.setValueTracker(vt, vt.Register(p.ID()))
	p.updateVtParams()
}

// unregisterPeer implements serverPeerSubscriber
func (v *vtSubscription) unregisterPeer(p *serverPeer) {
	vt := (*lpc.ValueTracker)(v)
	vt.Unregister(p.ID())
	p.setValueTracker(nil, nil)
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Etherbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("mining is not supported in light mode")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the mouse package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightEthereum) APIs() []rpc.API {
	apis := ethapi.GetAPIs(s.ApiBackend)
	apis = append(apis, s.engine.APIs(s.BlockChain().HeaderChain())...)
	return append(apis, []rpc.API{
		{
			Namespace: "mos",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "mos",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "mos",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "les",
			Version:   "1.0",
			Service:   NewPrivateLightAPI(&s.lesCommons),
			Public:    false,
		}, {
			Namespace: "lespay",
			Version:   "1.0",
			Service:   lpc.NewPrivateClientAPI(s.valueTracker),
			Public:    false,
		},
	}...)
}

func (s *LightEthereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightEthereum) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightEthereum) TxPool() *light.TxPool              { return s.txPool }
func (s *LightEthereum) Engine() consensus.Engine           { return s.engine }
func (s *LightEthereum) LesVersion() int                    { return int(ClientProtocolVersions[0]) }
func (s *LightEthereum) Downloader() *downloader.Downloader { return s.handler.downloader }
func (s *LightEthereum) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols returns all the currently configured network protocols to start.
func (s *LightEthereum) Protocols() []p2p.Protocol {
	return s.makeProtocols(ClientProtocolVersions, s.handler.runPeer, func(id enode.ID) interface{} {
		if p := s.peers.peer(id.String()); p != nil {
			return p.Info()
		}
		return nil
	}, s.dialCandidates)
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// light mouse protocol implementation.
func (s *LightEthereum) Start() error {
	log.Warn("Light client mode is an experimental feature")

	s.serverPool.start()
	// Start bloom request workers.
	s.wg.Add(bloomServiceThreads)
	s.startBloomHandlers(params.BloomBitsBlocksClient)
	s.handler.start()

	return nil
}

// Stop implements node.Lifecycle, terminating all internal goroutines used by the
// Mouse protocol.
func (s *LightEthereum) Stop() error {
	close(s.closeCh)
	s.serverPool.stop()
	s.valueTracker.Stop()
	s.peers.close()
	s.reqDist.close()
	s.odr.Stop()
	s.relay.Stop()
	s.bloomIndexer.Close()
	s.chtIndexer.Close()
	s.blockchain.Stop()
	s.handler.stop()
	s.txPool.Stop()
	s.engine.Close()
	s.pruner.close()
	s.eventMux.Stop()
	s.chainDb.Close()
	s.wg.Wait()
	log.Info("Light mouse stopped")
	return nil
}