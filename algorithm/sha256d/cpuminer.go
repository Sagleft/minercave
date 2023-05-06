package sha256d

import (
	"time"

	"minercave/net"
)

type HashRate struct {
	Rate       float64
	CPUMinerID int
}

type miningJob struct {
	Header []byte
	Offset int
}

type CPUMiner struct {
	validShares   uint
	staleShares   uint
	invalidShares uint
	devices       int
	startTime     uint
	HashRate      chan *HashRate
	miningJob     chan *miningJob
	pool          *net.Stratum
}

func NewCPUMiner(cfg *net.Config) (miner *CPUMiner) {
	stratum := net.StratumClient(cfg)
	miner = &CPUMiner{
		devices:   cfg.Threads,
		pool:      stratum,
		startTime: uint(time.Now().Unix()),
	}
	miner.miningJob = make(chan *miningJob, miner.devices)

	return
}

func (miner *CPUMiner) createJob() {
	miner.pool.Connect()

}

func (miner *CPUMiner) Mine() {
	go miner.createJob()

	dispatcher := NewDispatcher()
	dispatcher.Run()

	//Add Chunks on a regular basis
	for {
		//Get a new Chunk and split it accordingly to the machin settings
		for _, chunk := range NewChunkList(2, uint32(time.Now().Unix()), 2) {
			if len(dispatcher.ChunkQueueIn) < cap(dispatcher.ChunkQueueIn) {
				dispatcher.ChunkQueueIn <- chunk
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
