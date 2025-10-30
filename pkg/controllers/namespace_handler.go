package controllers

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func EnsureNamespace(ctx context.Context, c client.Client, namespaceName string) error {
	log := logrus.WithFields(logrus.Fields{
		"handler":   "namespace",
		"namespace": namespaceName,
	})

	if namespaceName == "" {
		err := fmt.Errorf("namespace name is not set in the configuration")
		log.Error(err)
		return err
	}

	namespace := &corev1.Namespace{}
	err := c.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Infof("Namespace %s does not exist, creating it", namespaceName)
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespaceName,
				},
			}
			if err := c.Create(ctx, namespace); err != nil {
				log.WithError(err).Error("unable to create namespace")
				return err
			}
			log.Infof("Successfully created namespace %s", namespaceName)
		} else {
			log.WithError(err).Error("unable to fetch namespace")
			return err
		}
	} else {
		log.Debugf("Namespace %s already exists", namespaceName)
	}

	return nil
}

func LabelNamespaceWithProjectID(ctx context.Context, c client.Client, namespaceName, projectID string) error {
	log := logrus.WithFields(logrus.Fields{
		"handler":   "namespace",
		"namespace": namespaceName,
		"projectID": projectID,
	})

	namespace := &corev1.Namespace{}
	if err := c.Get(ctx, types.NamespacedName{Name: namespaceName}, namespace); err != nil {
		log.WithError(err).Error("unable to fetch namespace for labeling")
		return err
	}
	log.Infof("Successfully fetched namespace: %s", namespace.Name)

	if namespace.Labels == nil {
		namespace.Labels = make(map[string]string)
	}

	if val, ok := namespace.Labels["rancher.k8s.binbash.org/project-name"]; ok && val == projectID {
		log.Infof("Namespace %s already has the correct project-name label", namespace.Name)
		return nil
	}

	namespace.Labels["rancher.k8s.binbash.org/project-name"] = projectID
	if err := c.Update(ctx, namespace); err != nil {
		log.WithError(err).Error("unable to update namespace with project-name label")
		return err
	}
	log.Infof("Successfully updated namespace %s with project-name label", namespace.Name)

	return nil
}
