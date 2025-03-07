package main

import (
	"flag"
	"os"

	"github.com/bryanpaget/statcan-namespace-auditor/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var dryRun bool

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Perform a dry run without deleting namespaces")
	flag.Parse()
	// Add Kubernetes core APIs to the scheme
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	// If you have any CRDs or additional schemes, add them here.
	// utilruntime.Must(yourcrd.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool

	// Command line flags
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	// Setup logging
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create a new manager to set up shared dependencies and start components
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "statcan-namespace-auditor",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Set up the NamespaceReconciler with the manager
	if err = (&controllers.NamespaceReconciler{
		Client: mgr.GetClient(),
		// Inject additional dependencies as needed, e.g., a logger or EntraID client.
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NamespaceReconciler")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
