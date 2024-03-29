/*
Copyright 2023 Compose, Borealis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	borealisdbiov1 "github.com/borealisdb/commons/borealisdb.io/v1"
	versioned "github.com/borealisdb/commons/generated/clientset/versioned"
	internalinterfaces "github.com/borealisdb/commons/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/borealisdb/commons/generated/listers/borealisdb.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BorealisClusterAccountInformer provides access to a shared informer and lister for
// BorealisClusterAccounts.
type BorealisClusterAccountInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.BorealisClusterAccountLister
}

type borealisClusterAccountInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBorealisClusterAccountInformer constructs a new informer for BorealisClusterAccount type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBorealisClusterAccountInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBorealisClusterAccountInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBorealisClusterAccountInformer constructs a new informer for BorealisClusterAccount type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBorealisClusterAccountInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BorealisdbV1().BorealisClusterAccounts(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.BorealisdbV1().BorealisClusterAccounts(namespace).Watch(context.TODO(), options)
			},
		},
		&borealisdbiov1.BorealisClusterAccount{},
		resyncPeriod,
		indexers,
	)
}

func (f *borealisClusterAccountInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBorealisClusterAccountInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *borealisClusterAccountInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&borealisdbiov1.BorealisClusterAccount{}, f.defaultInformer)
}

func (f *borealisClusterAccountInformer) Lister() v1.BorealisClusterAccountLister {
	return v1.NewBorealisClusterAccountLister(f.Informer().GetIndexer())
}
