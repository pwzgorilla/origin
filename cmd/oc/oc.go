package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/kubernetes/pkg/kubectl/scheme"

	"github.com/openshift/library-go/pkg/serviceability"
	"github.com/openshift/origin/pkg/oc/cli"
	"github.com/openshift/origin/pkg/version"

	// install all APIs
	apiinstall "github.com/openshift/origin/pkg/api/install"
	apilegacy "github.com/openshift/origin/pkg/api/legacy"
	_ "k8s.io/kubernetes/pkg/apis/autoscaling/install"
	_ "k8s.io/kubernetes/pkg/apis/batch/install"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	defer serviceability.BehaviorOnPanic(os.Getenv("OPENSHIFT_ON_PANIC"), version.Get())()
	defer serviceability.Profile(os.Getenv("OPENSHIFT_PROFILE")).Stop()

	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	fmt.Println("======", scheme.Scheme)
	fmt.Println("======", scheme.GroupFactoryRegistry)
	fmt.Println("======", scheme.Registry)
	apiinstall.InstallAll(scheme.Scheme, scheme.GroupFactoryRegistry, scheme.Registry)
	apilegacy.LegacyInstallAll(scheme.Scheme, scheme.Registry)

	basename := filepath.Base(os.Args[0])
	command := cli.CommandFor(basename)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
