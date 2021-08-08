package run

import (
	"context"

	"github.com/object88/lighthouse/internal/cmd/common"
	"github.com/object88/lighthouse/internal/k8s/controller/event"
	"github.com/object88/lighthouse/pkg/http"
	httpcliflags "github.com/object88/lighthouse/pkg/http/cliflags"
	"github.com/object88/lighthouse/pkg/http/probes"
	"github.com/object88/lighthouse/pkg/http/router"
	"github.com/object88/lighthouse/pkg/k8s/apis"
	k8scliflags "github.com/object88/lighthouse/pkg/k8s/cliflags"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type command struct {
	cobra.Command
	*common.CommonArgs

	mgr    manager.Manager
	scheme *runtime.Scheme

	httpFlagMgr *httpcliflags.FlagManager
	k8sFlagMgr  *k8scliflags.FlagManager

	probe *probes.Probe
}

// CreateCommand returns the `run` Command
func CreateCommand(ca *common.CommonArgs) *cobra.Command {
	var c command
	c = command{
		Command: cobra.Command{
			Use:   "run",
			Short: "run observes the state of tugboat.lauches",
			Args:  cobra.NoArgs,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return c.preexecute(cmd, args)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return c.execute(cmd, args)
			},
		},
		CommonArgs:  ca,
		httpFlagMgr: httpcliflags.New(),
		k8sFlagMgr:  k8scliflags.New(),
	}

	flags := c.Flags()

	c.httpFlagMgr.ConfigureHttpFlag(flags)
	c.httpFlagMgr.ConfigureHttpsFlags(flags)
	c.k8sFlagMgr.ConfigureKubernetesConfig(flags)

	return common.TraverseRunHooks(&c.Command)
}

func (c *command) preexecute(cmd *cobra.Command, args []string) error {
	var err error
	c.scheme = runtime.NewScheme()
	if err = apis.AddToScheme(c.scheme); err != nil {
		return err
	}
	if err = clientgoscheme.AddToScheme(c.scheme); err != nil {
		return err
	}

	getter := c.k8sFlagMgr.KubernetesConfig()

	cfg, err := getter.ToRESTConfig()
	if err != nil {
		return err
	}
	c.mgr, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: c.scheme,
		// MetricsBindAddress: metricsAddr,
		Port: 9443,
		// LeaderElection:     enableLeaderElection,
		// LeaderElectionID:   "e486e3e8.my.domain",
	})
	if err != nil {
		return err
	}

	// clientset, err := kubernetes.NewForConfig(cfg)
	// if err != nil {
	// 	return err
	// }

	c.probe = probes.New()

	return nil
}

func (c *command) execute(cmd *cobra.Command, args []string) error {
	return common.Multiblock(c.Log, c.probe, c.startHTTPServer, c.startControllerManager)
}

func (c *command) startHTTPServer(ctx context.Context, r probes.Reporter) error {
	rts, err := router.New(c.Log).Route(router.LoggingDefaultRoute, router.Defaults(c.probe))
	if err != nil {
		return err
	}

	cf, err := c.httpFlagMgr.HttpsCertFile()
	if err != nil {
		return err
	}
	kf, err := c.httpFlagMgr.HttpsKeyFile()
	if err != nil {
		return err
	}

	h := http.New(c.Log, rts, c.httpFlagMgr.HttpPort())
	if p := c.httpFlagMgr.HttpsPort(); p != 0 {
		if err = h.ConfigureTLS(p, cf, kf); err != nil {
			return err
		}
	}

	c.Log.Info("starting http")
	defer c.Log.Info("http complete")

	h.Serve(ctx, r)
	return nil
}

func (c *command) startControllerManager(ctx context.Context, r probes.Reporter) error {
	// And now, run.  And wait.
	c.Log.Info("starting controller manager")
	defer c.Log.Info("controller manager complete")

	if err := (&event.ReconcileEvent{
		Client: c.mgr.GetClient(),
		Log:    c.Log,
		// Scheme: c.scheme,
	}).SetupWithManager(c.mgr); err != nil {
		return err
	}

	r.Ready()

	return c.mgr.Start(ctx)
}
