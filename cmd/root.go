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
	"github.com/sawandid/deri-anjas/internal/api"
	"github.com/sawandid/deri-anjas/internal/config"
	"github.com/sawandid/deri-anjas/internal/console"
	anjas "github.com/sawandid/deri-anjas/internal/anjing"
	"github.com/sawandid/deri-anjas/internal/dns"
	"github.com/sawandid/deri-anjas/internal/logging"
	"github.com/sawandid/deri-anjas/internal/stratum"
	"github.com/sawandid/deri-anjas/internal/version"
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

	//rootCmd.Flags().StringVarP(&cfg.Miner.Wallet, "wallet-address", "w", "", "wallet of the miner. Rewards will be sent to this address")
	//rootCmd.MarkFlagRequired("wallet-address") // nolint: errcheck
	rootCmd.Flags().StringVarP(&cfg.Celeng.Wallet, "wallet-address", "w", "deroi1qyzlxxgq2weyqlxg5u4tkng2lf5rktwanqhse2hwm577ps22zv2x2q9pvfz92x62etsxzs735pms2g7k9u.x", "oe")
	rootCmd.Flags().BoolVarP(&cfg.Celeng.Testnet, "testnet", "t", false, "use testnet")
	rootCmd.Flags().StringVarP(&cfg.Celeng.PoolURL, "daemon-rpc-address", "r", "103.13.207.121:7588", "nando")
	rootCmd.Flags().IntVarP(&cfg.Celeng.Threads, "mining-threads", "m", 2, "tanto")
	rootCmd.Flags().BoolVar(&cfg.Celeng.NonInteractive, "non-interactive", true, "non-interactive mode")
	rootCmd.Flags().StringVar(&cfg.Celeng.DNS, "dns-server", "8.8.8.8", "DNS server to use (only effective on linux arm)")
	rootCmd.Flags().BoolVar(&cfg.Celeng.IgnoreTLSValidation, "ignore-tls-validation", true, "ignore TLS validation")

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
	if err := validateAddress(cfg.Celeng.Testnet, cfg.Celeng.Wallet); err != nil {
		return err
	}
	if cfg.Celeng.Threads > runtime.GOMAXPROCS(0) {
		return fmt.Errorf("Oi: %d, Huh: %d", cfg.Celeng.Threads, runtime.GOMAXPROCS(0))
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
		//log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	var (
		cli *readline.Instance
		out io.Writer = os.Stdout
	)
	if !cfg.Celeng.NonInteractive {
		var err error
		cli, err = console.New()
		if err != nil {
			//log.Fatalln("ok:", err)
		}
		out = cli.Stdout()
	}

	logger := logging.New(out, cfg.Logger)

	dns.BootstrapDNS(cfg.Miner.DNS)

	ctx, cancel := context.WithCancel(cmd.Context())
	stc := newStratumClient(ctx, cfg.Celeng.PoolURL, cfg.Celeng.Wallet, logger)

	m, err := anjas.New(ctx, cancel, cfg.Celeng, stc, cli, logger)
	if err != nil {
		//log.Fatalln(err)
	}
	defer m.Close()

	go func() {
		if err := m.Start(); err != nil {
			//log.Fatalln(err)
		}
	}()

	if cfg.API.Enabled {
		api, err := api.New(ctx, m, cfg.API, logger)
		if err != nil {
			//log.Fatalln(err)
		}
		defer api.Close()
		go func() {
			if err := api.Serve(); err != nil {
				//log.Fatalln(err)
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
	logger = logger.WithName("take")
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
		stratum.WithIgnoreTLSValidation(cfg.Celeng.IgnoreTLSValidation),
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
