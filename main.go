package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Guaderxx/cobra"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initLogger()
	initFlags()
}

var (
	iplist            = ""
	portList          = ""
	mode              = ""
	timeout     int8  = 0
	concurrency int32 = 0

	result = &sync.Map{}

	cmd = &cobra.Command{
		Use:     "scanner",
		Short:   "Tcp syn/connect port scanner",
		Long:    "xxx",
		Version: "2024/4/11",
		Run:     execute,
	}
)

func initFlags() {
	cmd.PersistentFlags().StringVarP(&iplist, "iplist", "i", "", "ip list")
	cmd.MarkPersistentFlagRequired("iplist")
	cmd.PersistentFlags().StringVarP(&portList, "port", "p", "22,23,53,80-139", "port list")
	cmd.PersistentFlags().StringVarP(&mode, "mode", "m", "syn", "scan mode")
	cmd.PersistentFlags().Int8VarP(&timeout, "timeout", "t", 3, "timeout")
	cmd.PersistentFlags().Int32VarP(&concurrency, "concurrency", "c", 5000, "concurrency")

}

func initLogger() {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// remove the directory from the source's filename
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}

		return a
	}

	tmp := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelError,
		ReplaceAttr: replace,
	}))
	slog.SetDefault(tmp)
}

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("execute error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

type IpPort struct {
	ip   string
	port int
}

func execute(cmd *cobra.Command, args []string) {
	if mode == "syn" {
		checkRoot()
	}
	ips, err := getIPList(iplist)
	if err != nil {
		slog.Error("parse flag ip-list error", slog.String("err", err.Error()))
		os.Exit(1)
	}
	ports, err := getPorts(portList)
	if err != nil {
		slog.Error("parse flag port error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Info("ips", slog.Any("ip-list", ips))
	slog.Info("ports", slog.Any("port-list", ports))

	tasks := []IpPort{}
	for _, ip := range ips {
		for _, port := range ports {
			tasks = append(tasks, IpPort{
				ip:   ip.String(),
				port: port,
			})
		}
	}

	RunTask(tasks)
	PrintResult()
}
