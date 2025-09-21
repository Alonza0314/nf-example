package service

import (
	"context"
	"io"
	"os"
	"runtime/debug"
	"sync"

	nf_context "github.com/Alonza0314/nf-example/internal/context"
	"github.com/Alonza0314/nf-example/internal/logger"
	"github.com/Alonza0314/nf-example/internal/sbi"
	"github.com/Alonza0314/nf-example/internal/sbi/processor"
	"github.com/Alonza0314/nf-example/pkg/app"
	"github.com/Alonza0314/nf-example/pkg/factory"
	"github.com/sirupsen/logrus"
)

type NfApp struct {
	cfg   *factory.Config
	nfCtx *nf_context.NFContext

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	sbiServer *sbi.Server
	processor *processor.Processor
}

var _ app.App = &NfApp{}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*NfApp, error) {
	// Initialize NF context from factory.NfConfig
	nf_context.InitNfContext()

	nf := &NfApp{
		cfg:   cfg,
		wg:    sync.WaitGroup{},
		nfCtx: nf_context.GetSelf(),
	}

	nf.SetLogEnable(cfg.GetLogEnable())
	nf.SetLogLevel(cfg.GetLogLevel())
	nf.SetReportCaller(cfg.GetLogReportCaller())

	nf.ctx, nf.cancel = context.WithCancel(ctx)

	// Initialize SBI server
	sbiServer := sbi.NewServer(nf, tlsKeyLogPath)
	nf.sbiServer = sbiServer

	processor, err := processor.NewProcessor(nf)
	if err != nil {
		return nf, err
	}
	nf.processor = processor

	return nf, nil
}

func (a *NfApp) Config() *factory.Config {
	return a.cfg
}

func (a *NfApp) Context() *nf_context.NFContext {
	return a.nfCtx
}

func (a *NfApp) Processor() *processor.Processor {
	return a.processor
}

func (a *NfApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}
	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *NfApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}
	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}
	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *NfApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}
	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (a *NfApp) Start() {
	defer func() {
		if p := recover(); p != nil {
			logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	a.sbiServer.Run(&a.wg)

	go a.listenShutdown(a.ctx)
	a.Wait()
}

func (a *NfApp) listenShutdown(ctx context.Context) {
	<-ctx.Done()
	a.terminateProcedure()
}

func (a *NfApp) Terminate() {
	a.cancel()
}

func (a *NfApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating ANYA...")
	a.sbiServer.Shutdown()
}

func (a *NfApp) Wait() {
	a.wg.Wait()
	logger.MainLog.Infof("ANYA terminated")
}
