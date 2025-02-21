package anjas

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chzyer/readline"
	"github.com/deroproject/derohe/astrobwt/astrobwtv3"
	"github.com/deroproject/derohe/block"
	"github.com/go-logr/logr"
	"github.com/jpillora/backoff"
	"github.com/teivah/broadcast"
	"github.com/sawandid/deri-anjas/internal/config"
	"github.com/sawandid/deri-anjas/internal/stratum"
)

var reportHashrateInterval = time.Second * 30

type Client struct {
	counter uint64 // Must be the first field. Otherwise atomic operations panic on arm7
	ctx     context.Context
	cancel  context.CancelFunc
	config  *config.Celeng
	stratum *stratum.Client
	console *readline.Instance
	logger  logr.Logger

	mu           sync.RWMutex
	job          *stratum.Job
	jobCounter   int64
	iterations   int
	hashrate     uint64
	mining       bool
	miningString string
	diffString   string
	heightString string

	shareCounter    uint64
	rejectedCounter uint64
}

func New(ctx context.Context, cancel context.CancelFunc, config *config.Celeng, stratum *stratum.Client, console *readline.Instance, logger logr.Logger) (*Client, error) {
	c := &Client{
		ctx:        ctx,
		cancel:     cancel,
		config:     config,
		stratum:    stratum,
		iterations: 100,
		console:    console,
	}
	c.setLogger(logger)
	return c, nil
}

func (c *Client) Close() error {
	if c.console != nil {
		return c.console.Close()
	}
	return nil
}

func (c *Client) Start() error {
	if c.config.Threads < 1 || c.iterations < 1 || c.config.Threads > 2048 {
		panic("Invalid parameters\n")
	}
	if c.config.Threads > 255 {
		c.logger.Error(nil, "okle.", "available", c.config.Threads)
		c.config.Threads = 255
	}

	go c.gatherStats()
	if c.config.NonInteractive {
		go c.noniSummary()
	}

	go c.getwork()

	for i := 0; i < c.config.Threads; i++ {
		go c.mineblock(i)
	}

	go c.reportHashrate()

	if !c.config.NonInteractive {
		c.startConsole()
	}
	return nil
}

func (c *Client) makeBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    time.Second,
		Max:    time.Second * 30,
		Factor: 1.5,
		Jitter: true,
	}
}

func (c *Client) getwork() {
	b := c.makeBackoff()
	rand.Seed(time.Now().UTC().UnixNano())
	for {
		if err := c.stratum.Dial(); err != nil {
			waitDuration := b.Duration()
			c.logger.Error(err, "oki", "server adress", c.config.PoolURL)
			c.logger.Info(fmt.Sprintf("de %f de", waitDuration.Seconds()))
			time.Sleep(waitDuration)
			continue
		}

		jobListener := c.stratum.NewJobListener(4)
		defer jobListener.Close()

		respListener := c.stratum.NewResponseListener(2)
		go c.listenStratumResponses(respListener)

		for {
			select {
			case j := <-jobListener.Ch():
				c.mu.Lock()
				c.job = j
				c.jobCounter++
				c.mu.Unlock()
			case <-c.ctx.Done():
				return
			}
		}
	}
}

func (c *Client) listenStratumResponses(l *broadcast.Listener[*stratum.Response]) {
	defer l.Close()
	for range l.Ch() {
		c.shareCounter = uint64(c.stratum.GetTotalShares())
		c.rejectedCounter = uint64(c.stratum.GetTotalShares() - c.stratum.GetAcceptedShares())
	}
}

func (c *Client) mineblock(tid int) {
	var diff big.Int
	var work [block.MINIBLOCK_SIZE]byte

	var randomBuf [12]byte

	rand.Read(randomBuf[:]) //#nosec G404

	time.Sleep(time.Millisecond * 500)

	nonceBuf := work[block.MINIBLOCK_SIZE-5:] //since slices are linked, it modifies parent
	runtime.LockOSThread()
	threadaffinity()

	var localJobCounter int64

	i := uint32(0)

	for {
		c.mu.RLock()
		myjob := c.job
		localJobCounter = c.jobCounter
		c.mu.RUnlock()
		if myjob == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		n, err := hex.Decode(work[:], []byte(myjob.Blob))
		if err != nil || n != block.MINIBLOCK_SIZE {
			c.logger.Error(err, "unijei", "blockwork", myjob.Blob, "n", n, "job", myjob)
			time.Sleep(time.Millisecond * 500)
			continue
		}

		copy(work[block.MINIBLOCK_SIZE-12:], randomBuf[:]) // add more randomization in the mix
		work[block.MINIBLOCK_SIZE-1] = byte(tid)
		diff.SetString(strconv.Itoa(int(myjob.Difficulty)), 10)

		if work[0]&0xf != 1 { // check  version
			c.logger.Error(nil, "cahala", "version", work[0]&0x1f)
			time.Sleep(time.Millisecond * 500)
			continue
		}

		for localJobCounter == c.jobCounter { // update job when it comes, expected rate 2 per second
			if !c.stratum.IsConnected() {
				time.Sleep(time.Millisecond * 500)
				continue
			}
			i++
			binary.BigEndian.PutUint32(nonceBuf, i)

			powhash := astrobwtv3.AstroBWTv3(work[:])
			atomic.AddUint64(&c.counter, 1)

			if CheckPowHashBig(powhash, &diff) { // note we are doing a local, NW might have moved meanwhile
				c.logger.V(1).Info("oyiers", "difficulty", myjob.Difficulty, "height", myjob.Height)
				func() {
					defer c.recover(1) // nolint: errcheck
					nonce := work[len(work)-12:]
					share := stratum.NewShare(myjob.ID, fmt.Sprintf("%x", nonce), fmt.Sprintf("%x", powhash[:]))
					if err := c.stratum.SubmitShare(share); err != nil {
						c.logger.Error(err, "raiso")
					}
				}()
			}
		}
	}
}

func (c *Client) recover(level int) (err error) {
	if r := recover(); r != nil {
		err = fmt.Errorf("Recovered r:%+v stack %s", r, string(debug.Stack()))
		c.logger.V(level).Error(nil, "Recovered ", "error", r, "stack", string(debug.Stack()))
	}
	return
}

func (c *Client) reportHashrate() {
	ticker := time.NewTicker(reportHashrateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.stratum.ReportHashrate(stratum.NewReport(c.GetHashrate())); err != nil {
				c.logger.Error(err, "raizo")
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) GetHashrate() uint64 {
	return c.hashrate
}

func (c *Client) GetTotalShares() uint64 {
	return c.shareCounter
}

func (c *Client) GetAcceptedShares() uint64 {
	return c.shareCounter - c.rejectedCounter
}

func (c *Client) GetRejectedShares() uint64 {
	return c.rejectedCounter
}

func (c *Client) GetPoolURL() string {
	return c.config.PoolURL
}
