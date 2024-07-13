package autostager

import (
	"github.com/synoti21/auto-stager/internal/driver"
	autostager "github.com/synoti21/auto-stager/internal/driver"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Manager struct {
	AutostagerClient *autostager.AutostagerClient
}

func NewManager(kube client.Client, scheme *runtime.Scheme) (*Manager, error) {
	autostagerClient, err := driver.NewAutostagerClient(kube, scheme)
	if err != nil {
		return nil, err
	}
	return &Manager{
		AutostagerClient: autostagerClient,
	}, nil
}
