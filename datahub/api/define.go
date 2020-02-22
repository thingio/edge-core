package api

import "github.com/thingio/edge-core/common/proto/resource"

type ResourceWatcher func(*resource.Resource)

type DatahubApi interface {
	/* GET api: non-deleted resource by id */
	GetResource(kind resource.Kind, id string) (*resource.Resource, error)

	/* LIST api: non-deleted resources */
	ListResources(kind resource.Kind) ([]*resource.Resource, error)

	/* WATCH api: include the ones editor itself made */
	WatchResource(kind resource.Kind, watcher ResourceWatcher) error

	/* UPDATE api: save the entire resource */
	SaveResource(r *resource.Resource) error

	/* DELETE api: soft deletion, which marks version of the resource to 0 */
	DeleteResource(kind resource.Kind, id string) error
}
