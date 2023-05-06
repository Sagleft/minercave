/*
Author: Mathieu Mailhos
Filename: mining.go
Description: Functions for mining a Block Header
*/

package sha256d

import (
	"log"
	"minercave/algorithm/sha256d/block"
	"minercave/utils"
	"strconv"
	"time"
)

// MaxNonce standard value
const MaxNonce uint32 = 4294967295

// HashCountSpan counter big enough to avoid mutex bottleneck
const HashCountSpan uint32 = 200000

// Miner entity defined by an ID. Worker.
type Miner struct {
	ID              int
	MiningPool      chan chan Chunk
	BlockChannelIn  chan Chunk
	BlockChannelOut chan Chunk
	quit            chan bool
}

// NewMiner Creates Miner 'Worker'
func NewMiner(id int, miningpool chan chan Chunk, outchan chan Chunk) Miner {
	log.Printf("New Miner created.")
	return Miner{
		ID:              id,
		MiningPool:      miningpool,
		BlockChannelIn:  make(chan Chunk),
		BlockChannelOut: outchan,
		quit:            make(chan bool)}
}

// Start mining: receive block channels and execute them
func (mine Miner) Start() {
	go func() {
		for {
			//We register the mine into the mining pool
			mine.MiningPool <- mine.BlockChannelIn
			log.Printf("Miner " + strconv.Itoa(mine.ID) + " available.")
			select {
			//We then receive a chunk to work on or we quit
			case job := <-mine.BlockChannelIn:
				log.Printf("Miner " + strconv.Itoa(mine.ID) + " starts mining.")
				success, chunk := mine.mining(job)
				if success {
					//Send Back to dispatcher for validation, to be sent back to Websocket
					chunk.Valid = true
					utils.LOG_SUCCESS.Println("Verified Chunk")
					mine.BlockChannelOut <- chunk
				}
			case <-mine.quit:
				return
			}
		}
	}()
}

// Stop tells the Miner to stop working
func (mine Miner) Stop() {
	go func() {
		log.Printf("Mine " + strconv.Itoa(mine.ID) + " stopped.")
		mine.quit <- true
	}()
}

// Mining a blockheader and returning the chunk including proper nonce value if suceeded. Splited into two to avoid useless checks and increment if the loging is not activated. Mining can not take more time than 1 second as the block header expires due to the epoch time field changing constantly.
func (mine *Miner) mining(chunk Chunk) (bool, Chunk) {
	// var hashBig big.Int
	// var resultDiff big.Int
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	// if !log.Activated {
	// 	for nonce := chunk.StartNonce; nonce < chunk.EndNonce; nonce++ {
	// 		select {
	// 		case <-timeout:
	// 			//Timeout
	// 			return false, chunk
	// 		default:
	// 			//Success
	// 			chunk.Block.Nonce = nonce
	// 			if hash := block.Doublesha256BlockHeader(chunk.Block); hash < chunk.Target {
	// 				return true, chunk
	// 			}
	// 		}
	// 	}
	// } else {

	utils.LOG_INFO.Println(" ==> Target:" + chunk.Target)
	for count, nonce := uint32(0), chunk.StartNonce; nonce < chunk.EndNonce; nonce, count = nonce+1, count+1 {
		select {
		case <-timeout:
			//Timeout
			log.Printf("Timeout, moving to next block. " + strconv.Itoa(int(count)) + " operations done on this block by Miner " + strconv.Itoa(mine.ID))
			return false, chunk
		default:
			//Success
			chunk.Block.Nonce = nonce
			hash := block.Doublesha256BlockHeader(chunk.Block)
			// byteHash := block.Doublesha256BlockHeaderbyByte(chunk.Block)

			// utils.HashToBig(byteHash, &hashBig)
			// poolMatch := hashBig.Cmp(big.NewInt(1)) <= 0
			poolMatch2 := hash < chunk.Target
			// target := big.NewInt(2)

			if poolMatch2 {

				// utils.CalculateDifficulty(target, &resultDiff)
				// diff := utils.Difficulty(resultDiff.Uint64())

				// log.Printf(" %s %s \n", hashBig.String(), resultDiff.String())
				// log.Printf(" %s %s \n", utils.BigToHashString(big.NewInt(1)), diff.String())

				// utils.LOG_WARN(" ===> %v %v %s %s \n", poolMatch, poolMatch2, hash, chunk.Target)
				// log.IncrementBlockCount()
				log.Printf("NEW BLOCK FOUND!! Nonce:" + strconv.Itoa(int(nonce)) + " Miner:" + strconv.Itoa(mine.ID) + " Hash:" + hash)
				return true, chunk
			}
			// utils.LOG_INFO.Println(" ==> Nonce:" + strconv.Itoa(int(nonce)) + " Miner:" + strconv.Itoa(mine.ID) + " Hash:" + hash)
			// if count == HashCountSpan {
			// 	// log.IncrementHashCount(count)
			// 	count = 0
			// }
		}
	}
	// }
	//MaxNonce Reached
	return false, chunk
}
