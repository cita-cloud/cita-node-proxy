/*
 * Copyright Rivtower Technologies LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	citacloudv1 "github.com/cita-cloud/cita-node-operator/api/v1"
	"github.com/cita-cloud/cita-node-proxy/pkg"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var K8sClient client.Client

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

func main() {
	router := gin.Default()
	router.GET("/ping", pong)

	router.GET("/backups/:namespace/:name", getBackupByNamespaceName)
	router.POST("/backups", postBackup)
	router.DELETE("/backups/:namespace/:name", deleteBackup)

	router.GET("/restores/:namespace/:name", getRestoreByNamespaceName)
	router.POST("/restores", postRestore)
	router.DELETE("/restores/:namespace/:name", deleteRestore)

	K8sClient, _ = pkg.InitK8sClient()

	router.Run()
}

type Backup struct {
	Kind       string                   `json:"kind"`
	ApiVersion string                   `json:"apiVersion"`
	Metadata   ObjectMeta               `json:"metadata"`
	Spec       citacloudv1.BackupSpec   `json:"spec,omitempty"`
	Status     citacloudv1.BackupStatus `json:"status,omitempty"`
}

type ObjectMeta struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": "pong"})
}

func getBackupByNamespaceName(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	k8sBackup := &citacloudv1.Backup{}
	err := K8sClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, k8sBackup)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}

	var backup Backup

	backup.Kind = k8sBackup.Kind
	backup.ApiVersion = k8sBackup.APIVersion

	backup.Metadata.Name = k8sBackup.Name
	backup.Metadata.Namespace = k8sBackup.Namespace

	backup.Spec.Chain = k8sBackup.Spec.Chain
	backup.Spec.Namespace = k8sBackup.Spec.Namespace
	backup.Spec.Node = k8sBackup.Spec.Node
	backup.Spec.DeployMethod = k8sBackup.Spec.DeployMethod
	backup.Spec.StorageClass = k8sBackup.Spec.StorageClass
	backup.Spec.Action = k8sBackup.Spec.Action
	backup.Spec.Image = k8sBackup.Spec.Image
	backup.Spec.PullPolicy = k8sBackup.Spec.PullPolicy

	backup.Status.Status = k8sBackup.Status.Status
	backup.Status.EndTime = k8sBackup.Status.EndTime
	backup.Status.StartTime = k8sBackup.Status.StartTime
	backup.Status.Actual = k8sBackup.Status.Actual

	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": backup})
}

func postBackup(c *gin.Context) {
	var backup Backup
	if err := c.BindJSON(&backup); err != nil {
		return
	}
	k8sBackup := &citacloudv1.Backup{}

	k8sBackup.Name = backup.Metadata.Name
	k8sBackup.Namespace = backup.Metadata.Namespace

	k8sBackup.Spec.Chain = backup.Spec.Chain
	k8sBackup.Spec.Namespace = backup.Spec.Namespace
	k8sBackup.Spec.Node = backup.Spec.Node
	k8sBackup.Spec.DeployMethod = backup.Spec.DeployMethod
	k8sBackup.Spec.StorageClass = backup.Spec.StorageClass
	k8sBackup.Spec.Action = backup.Spec.Action
	k8sBackup.Spec.Image = backup.Spec.Image
	k8sBackup.Spec.PullPolicy = backup.Spec.PullPolicy

	err := K8sClient.Create(context.Background(), k8sBackup)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	// set active status
	backup.Status.Status = citacloudv1.JobActive
	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": backup})
}

func deleteBackup(c *gin.Context) {
	ctx := context.Background()

	namespace := c.Param("namespace")
	name := c.Param("name")

	backup := &citacloudv1.Backup{}
	err := K8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, backup)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	err = K8sClient.Delete(ctx, backup, client.GracePeriodSeconds(0))
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": nil})
}

type Restore struct {
	Kind       string                    `json:"kind"`
	ApiVersion string                    `json:"apiVersion"`
	Metadata   ObjectMeta                `json:"metadata"`
	Spec       citacloudv1.RestoreSpec   `json:"spec,omitempty"`
	Status     citacloudv1.RestoreStatus `json:"status,omitempty"`
}

func getRestoreByNamespaceName(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	k8sRestore := &citacloudv1.Restore{}
	err := K8sClient.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, k8sRestore)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}

	var restore Restore

	restore.Kind = k8sRestore.Kind
	restore.ApiVersion = k8sRestore.APIVersion

	restore.Metadata.Name = k8sRestore.Name
	restore.Metadata.Namespace = k8sRestore.Namespace

	restore.Spec.Chain = k8sRestore.Spec.Chain
	restore.Spec.Namespace = k8sRestore.Spec.Namespace
	restore.Spec.Node = k8sRestore.Spec.Node
	restore.Spec.DeployMethod = k8sRestore.Spec.DeployMethod
	restore.Spec.Backup = k8sRestore.Spec.Backup
	restore.Spec.Action = k8sRestore.Spec.Action
	restore.Spec.Image = k8sRestore.Spec.Image
	restore.Spec.PullPolicy = k8sRestore.Spec.PullPolicy

	restore.Status.Status = k8sRestore.Status.Status
	restore.Status.EndTime = k8sRestore.Status.EndTime
	restore.Status.StartTime = k8sRestore.Status.StartTime

	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": restore})
}

func postRestore(c *gin.Context) {
	var restore Restore
	if err := c.BindJSON(&restore); err != nil {
		return
	}

	k8sRestore := &citacloudv1.Restore{}

	k8sRestore.Name = restore.Metadata.Name
	k8sRestore.Namespace = restore.Metadata.Namespace

	k8sRestore.Spec.Chain = restore.Spec.Chain
	k8sRestore.Spec.Namespace = restore.Spec.Namespace
	k8sRestore.Spec.Node = restore.Spec.Node
	k8sRestore.Spec.DeployMethod = restore.Spec.DeployMethod
	k8sRestore.Spec.Backup = restore.Spec.Backup
	k8sRestore.Spec.Action = restore.Spec.Action
	k8sRestore.Spec.Image = restore.Spec.Image
	k8sRestore.Spec.PullPolicy = restore.Spec.PullPolicy

	err := K8sClient.Create(context.Background(), k8sRestore)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	// set active status
	restore.Status.Status = citacloudv1.JobActive
	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": restore})
}

func deleteRestore(c *gin.Context) {
	ctx := context.Background()

	namespace := c.Param("namespace")
	name := c.Param("name")

	restore := &citacloudv1.Restore{}
	err := K8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, restore)
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	err = K8sClient.Delete(ctx, restore, client.GracePeriodSeconds(0))
	if err != nil {
		c.JSON(500, gin.H{"msg": err, "status": FAIL, "data": nil})
	}
	c.JSON(200, gin.H{"msg": nil, "status": SUCCESS, "data": nil})
}
