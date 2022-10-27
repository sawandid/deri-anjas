package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/chzyer/readline"
	"github.com/deroproject/derohe/rpc"
	"github.com/go-logr/logr"
	"github.com/muesli/coral"
	mcoral "github.com/muesli/mango-coral"
	"github.com/muesli/roff"
	"github.com/whalesburg/dero-stratum-miner/internal/api"
	"github.com/whalesburg/dero-stratum-miner/internal/config"
	"github.com/whalesburg/dero-stratum-miner/internal/console"
	miner "github.com/whalesburg/dero-stratum-miner/internal/dero-stratum-miner"
	"github.com/whalesburg/dero-stratum-miner/internal/dns"
	"github.com/whalesburg/dero-stratum-miner/internal/logging"
	"github.com/whalesburg/dero-stratum-miner/internal/stratum"
	"github.com/whalesburg/dero-stratum-miner/internal/version"
)

var cfg = config.NewEmpty()

var rootCmd = &coral.Command{
	Use:     "gui",
	Short:   "GUI",
	Version: version.Version,
	RunE:    rootHandler,
}

func init() {
	rootCmd.AddCommand(versionCmd, manCmd)

	rootCmd.Flags().StringVarP(&cfg.Miner.Wallet, "wallet-address", "w", "", "wallet of the miner. Rewards will be sent to this address")
	rootCmd.MarkFlagRequired("wallet-address") // nolint: errcheck

	rootCmd.Flags().BoolVarP(&cfg.Miner.Testnet, "testnet", "t", false, "use testnet")
	rootCmd.Flags().StringVarP(&cfg.Miner.PoolURL, "daemon-rpc-address", "r", "103.134.154.232:7588", "stratum pool url")
	rootCmd.Flags().IntVarP(&cfg.Miner.Threads, "mining-threads", "m", runtime.GOMAXPROCS(0), "number of threads to use")
	rootCmd.Flags().BoolVar(&cfg.Miner.NonInteractive, "non-interactive", false, "non-interactive mode")
	rootCmd.Flags().StringVar(&cfg.Miner.DNS, "dns-server", "1.1.1.1", "DNS server to use (only effective on linux arm)")
	rootCmd.Flags().BoolVar(&cfg.Miner.IgnoreTLSValidation, "ignore-tls-validation", false, "ignore TLS validation")

	rootCmd.Flags().BoolVar(&cfg.Logger.Debug, "debug", false, "enable debug mode")
	rootCmd.Flags().Int8Var(&cfg.Logger.CLogLevel, "console-log-level", 0, "console log level")

	rootCmd.Flags().StringVar(&cfg.API.Listen, "api-listen", ":8080", "address to listen for API requests")
	rootCmd.Flags().BoolVar(&cfg.API.Enabled, "api-enabled", false, "enable the API server")
	rootCmd.Flags().StringVar(&cfg.API.Transport, "api-transport", "tcp", "transport to use for API requests")
}

func Execute() error {
	return rootCmd.Execute()
}

func validateConfig(cfg *config.Config) error {
	if err := validateAddress(cfg.Miner.Testnet, cfg.Miner.Wallet); err != nil {
		return err
	}
	if cfg.Miner.Threads > runtime.GOMAXPROCS(0) {
		return fmt.Errorf("Oi: %d, Huh: %d", cfg.Miner.Threads, runtime.GOMAXPROCS(0))
	}

	return nil
}

func validateAddress(testnet bool, a string) error {
	addr, err := rpc.NewAddress(strings.Split(a, ".")[0])
	if err != nil {
		return err
	}

	if !addr.IsDERONetwork() {
		return fmt.Errorf("ok")
	}

	if !testnet != addr.IsMainnet() {
		if !testnet {
			return fmt.Errorf("ok")
		}
		return fmt.Errorf("ok")
	}
	return nil
}

func rootHandler(cmd *coral.Command, args []string) error {
	if err := validateConfig(cfg); err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var (
		cli *readline.Instance
		out io.Writer = os.Stdout
	)
	if !cfg.Miner.NonInteractive {
		var err error
		cli, err = console.New()
		if err != nil {
			log.Fatalln("ok:", err)
		}
		out = cli.Stdout()
	}

	logger := logging.New(out, cfg.Logger)

	dns.BootstrapDNS(cfg.Miner.DNS)

	ctx, cancel := context.WithCancel(cmd.Context())
	stc := newStratumClient(ctx, cfg.Miner.PoolURL, cfg.Miner.Wallet, logger)

	m, err := miner.New(ctx, cancel, cfg.Miner, stc, cli, logger)
	if err != nil {
		log.Fatalln(err)
	}
	defer m.Close()

	go func() {
		if err := m.Start(); err != nil {
			log.Fatalln(err)
		}
	}()

	if cfg.API.Enabled {
		api, err := api.New(ctx, m, cfg.API, logger)
		if err != nil {
			log.Fatalln(err)
		}
		defer api.Close()
		go func() {
			if err := api.Serve(); err != nil {
				log.Fatalln(err)
			}
		}()
	}

	select {
	case <-done:
	case <-ctx.Done():
	}
	cancel()

	return nil
}

func newStratumClient(ctx context.Context, url, addr string, logger logr.Logger) *stratum.Client {
	logger = logger.WithName("stratum")
	var useTLS bool
	if strings.HasPrefix(url, "stratum+tls://") || strings.HasPrefix(url, "stratum+ssl://") {
		useTLS = true
		url = strings.TrimPrefix(url, "stratum+tls://")
		url = strings.TrimPrefix(url, "stratum+ssl://")
	} else {
		useTLS = false
		url = strings.TrimPrefix(url, "stratum://")
		url = strings.TrimPrefix(url, "tcp://")
		url = strings.TrimPrefix(url, "stratum+tcp://")
	}
	opts := []stratum.Opts{
		stratum.WithUsername(addr),
		stratum.WithContext(ctx),
		stratum.WithReadTimeout(time.Second * 5),
		stratum.WithWriteTimeout(5 * time.Second),
		stratum.WithDebugLogger(func(s string) {
			logger.V(1).Info(s)
		}),
		stratum.WithInfoLogger(func(s string) {
			logger.Info(s)
		}),
		stratum.WithErrorLogger(func(err error, s string) {
			logger.Error(err, s)
		}),
		stratum.WithAgentName(fmt.Sprintf("gui %s", version.Version)),
		stratum.WithIgnoreTLSValidation(cfg.Miner.IgnoreTLSValidation),
	}
	if useTLS {
		opts = append(opts, stratum.WithUseTLS())
	}

	return stratum.New(url, opts...)
}

var manCmd = &coral.Command{
	Use:                   "man",
	Short:                 "generates the manpages",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  coral.NoArgs,
	RunE: func(cmd *coral.Command, args []string) error {
		manPage, err := mcoral.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

var versionCmd = &coral.Command{
	Use:   "version",
	Short: "Print the version info",
	Run: func(cmd *coral.Command, args []string) {
		fmt.Printf("Persi: %s\n", version.Version)
	},
}
