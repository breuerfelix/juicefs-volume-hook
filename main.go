package main

import (
	"flag"
	"os"
	"strings"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	var probeAddr string
	var annotation bool
	var storageClasses string
	var certDir, keyName, certName string

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&annotation, "pod-annotation", false, "Only change mounts for pods with a given annotation.")
	flag.StringVar(&storageClasses, "storage-classes", "", "Only change mounts for a given storageClassName.")

	flag.StringVar(&certDir, "cert-dir", "", "Folder where key-name and cert-name are located.")
	flag.StringVar(&keyName, "key-name", "", "Filename for .key file.")
	flag.StringVar(&certName, "cert-name", "", "Filename for .cert file.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	storageClassList := []string{}
	if storageClasses != "" {
		storageClassList = strings.Split(storageClasses, ",")
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         false,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Server uses default values if provided paths are empty
	server := &webhook.Server{
		CertDir:  certDir,
		KeyName:  keyName,
		CertName: certName,
	}

	server.Register("/mutate", &webhook.Admission{Handler: &podWebhook{
		Client:         mgr.GetClient(),
		Annotation:     annotation,
		StorageClasses: storageClassList,
	}})

	mgr.Add(server)

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
