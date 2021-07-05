package collectors

import (
	"github.com/openshift/ocs-operator/metrics/internal/options"
	"github.com/prometheus/client_golang/prometheus"
	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
	cephv1listers "github.com/rook/rook/pkg/client/listers/ceph.rook.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

const (
	// name of the project/exporter
	ns = "ocs"
	// component within the project/exporter
	cephcluster = "cephcluster"
)

var _ prometheus.Collector = &CephClusterCollector{}

// CephClusterCollector is a custom collector for CephCluster Custom Resource
type CephClusterCollector struct {
	ClusterState      *prometheus.Desc
	Informer          cache.SharedIndexInformer
	AllowedNamespaces []string
}

// NewCephClusterCollector constructs a collector
func NewCephClusterCollector(opts *options.Options) *CephClusterCollector {
	client, err := rookclient.NewForConfig(opts.Kubeconfig)
	if err != nil {
		klog.Error(err)
	}

	lw := cache.NewListWatchFromClient(client.CephV1().RESTClient(), "CephClusters", metav1.NamespaceAll, fields.Everything())
	sharedIndexInformer := cache.NewSharedIndexInformer(lw, &cephv1.CephCluster{}, 0, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})

	return &CephClusterCollector{
		ClusterState: prometheus.NewDesc(
			prometheus.BuildFQName(ns, cephcluster, "state"),
			`Status of CephCluster. 1=Present`,
			[]string{"name", "namespace", "status"},
			nil,
		),
		Informer:          sharedIndexInformer,
		AllowedNamespaces: opts.AllowedNamespaces,
	}
}

// Run starts CephCluster informer
func (c *CephClusterCollector) Run(stopCh <-chan struct{}) {
	go c.Informer.Run(stopCh)
}

// Describe implements prometheus.Collector interface
func (c *CephClusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{}

	for _, d := range ds {
		ch <- d
	}
}

// Collect implements prometheus.Collector interface
func (c *CephClusterCollector) Collect(ch chan<- prometheus.Metric) {
	CephClusterLister := cephv1listers.NewCephClusterLister(c.Informer.GetIndexer())
	CephClusters := getAllCephClusters(CephClusterLister, c.AllowedNamespaces)

	if len(CephClusters) > 0 {
		c.collectCephClusters(CephClusters, ch)
	}
}

func getAllCephClusters(lister cephv1listers.CephClusterLister, namespaces []string) (CephClusters []*cephv1.CephCluster) {
	var tempCephClusters []*cephv1.CephCluster
	var err error
	if len(namespaces) == 0 {
		CephClusters, err = lister.List(labels.Everything())
		if err != nil {
			klog.Errorf("couldn't list CephClusters. %v", err)
		}
		return
	}
	for _, namespace := range namespaces {
		tempCephClusters, err = lister.CephClusters(namespace).List(labels.Everything())
		if err != nil {
			klog.Errorf("couldn't list CephClusters in namespace %s. %v", namespace, err)
			continue
		}
		CephClusters = append(CephClusters, tempCephClusters...)
	}
	return
}

func (c *CephClusterCollector) collectCephClusters(CephClusters []*cephv1.CephCluster, ch chan<- prometheus.Metric) {
	for _, CephCluster := range CephClusters {
		ch <- prometheus.MustNewConstMetric(c.ClusterState,
			prometheus.GaugeValue, 1,
			CephCluster.Name,
			CephCluster.Namespace)
	}
}
