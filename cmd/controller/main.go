package main

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/connectioninterface"
	"github.com/octopipe/cloudx/internal/controller/runner"
	"github.com/octopipe/cloudx/internal/controller/sharedinfra"
	"github.com/octopipe/cloudx/internal/provider"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     "0",
		Port:                   9443,
		HealthProbeBindAddress: ":8001",
		LeaderElection:         false,
		LeaderElectionID:       "dec90f54.circlerr.io",
	})
	if err != nil {
		panic(err)
	}

	k8sClient := kubernetes.NewForConfigOrDie(mgr.GetConfig())

	provider := provider.NewProvider(mgr.GetClient())
	sharedInfraController := sharedinfra.NewController(
		logger,
		mgr.GetClient(),
		mgr.GetScheme(),
		provider,
	)

	runnerController := runner.NewController(
		logger,
		mgr.GetClient(),
		mgr.GetScheme(),
		k8sClient,
	)

	if err := sharedInfraController.SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := runnerController.SetupWithManager(mgr); err != nil {
		panic(err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		panic(err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		panic(err)
	}

	connectionInterfaceRepository := connectioninterface.NewK8sRepository(mgr.GetClient())

	sharedInfraRPCServer := sharedinfra.NewRPCServer(mgr.GetClient(), logger)
	connectionInterfaceRPCServer := connectioninterface.NewConnectionInterfaceRPCHandler(connectionInterfaceRepository)
	rpc.Register(sharedInfraRPCServer)
	rpc.Register(connectionInterfaceRPCServer)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":9000")
	if e != nil {
		panic(err)
	}

	logger.Info("start rpc server")
	go http.Serve(l, nil)

	logger.Info("start controllers")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		panic(err)
	}
}
