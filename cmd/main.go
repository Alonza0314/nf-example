package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/Alonza0314/nf-example/pkg/service"
	"github.com/urfave/cli"

	logger_util "github.com/free5gc/util/logger"
	"github.com/free5gc/util/version"
)

var NF *service.NfApp

func main() {
	// Recover from panic and print stack to log
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()
	// command line options
	app := cli.NewApp()
	app.Name = "anya"
	app.Usage = "SPYxFamily"
	app.Action = action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringSliceFlag{
			Name:  "log, l",
			Usage: "Output NF log to `FILE`",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("ANYA Run Error: %v\n", err)
	}
}

// called by app.Run automatically if cli's children is not executed
func action(cliCtx *cli.Context) error {
	// ./bin/nf -c config/nfcfg.yaml
	// StringSlice is empty
	tlsKeyLogPath, err := initLogFile(cliCtx.StringSlice("log"))
	if err != nil {
		return err
	}
	logger.MainLog.Infoln("Anya version: ", version.GetVersion())
	// String is config/nfcfg.yaml
	cfg, err := factory.ReadConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}
	// set cfg to global variable in factory package
	factory.NfConfig = cfg

	// handle SIGINT and SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh  // Wait for interrupt signal to gracefully shutdown
		cancel() // Notify each goroutine and wait them stopped
	}()
	// ctx: cancel when SIGINT or SIGTERM received
	// cfg: configuration read from config file
	// tlsKeyLogPath: path to store TLS key log for wireshark
	nf, err := service.NewApp(ctx, cfg, tlsKeyLogPath)
	if err != nil {
		return err
	}
	NF = nf

	nf.Start()

	return nil
}

func initLogFile(logNfPath []string) (string, error) {
	logTlsKeyPath := ""

	for _, path := range logNfPath {
		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
			return "", err
		}

		if logTlsKeyPath != "" {
			continue
		}
		nfDir, _ := filepath.Split(path)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return "", err
		}
		_, name := filepath.Split(factory.NfDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}
	return logTlsKeyPath, nil
}
