package prebind

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// StatelessPreBindExample is an example of a simple plugin that has no state
// and implements only one hook for prebind.
type StatelessPreBindExample struct{}

var _ framework.PreBindPlugin = StatelessPreBindExample{}

// Name is the name of the plugin used in Registry and configurations.
const Name = "stateless-prebind-plugin-example"

// Name returns name of the plugin. It is used in logs, etc.
func (sr StatelessPreBindExample) Name() string {
	return Name
}

// PreBind is the functions invoked by the framework at "prebind" extension point.
func (sr StatelessPreBindExample) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	if pod.Namespace != "foo" {
		return framework.NewStatus(framework.Unschedulable, "only pods from 'foo' namespace are allowed")
	}
	return nil
}

// New initializes a new plugin and returns it.
func New(_ context.Context, _ *runtime.Unknown, _ framework.Handle) (framework.Plugin, error) {
	return &StatelessPreBindExample{}, nil
}
