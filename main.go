package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func main() {
	// Инициализация Kubernetes клиента с использованием настроек из пода
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Маршруты
	r.GET("/", indexHandler)
	r.GET("/pods", podsHandler)
	r.POST("/deploy/:name", deployHandler)
	r.POST("/delete/:name", deleteHandler)

	r.Run(":8080")
}

func indexHandler(c *gin.Context) {
	pods, err := getPods()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"pods": nil, "error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"pods": pods, "error": nil})
}

func podsHandler(c *gin.Context) {
	pods, err := getPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "pods": pods})
}

func deployHandler(c *gin.Context) {
	name := c.Param("name")
	err := deploy(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Deployed " + name})
}

func deleteHandler(c *gin.Context) {
	name := c.Param("name")
	err := delete(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Deleted " + name})
}

func getPods() ([]string, error) {
	pods, err := clientset.CoreV1().Pods("axis-testing").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	return podNames, nil
}

func deploy(name string) error {
	// Логика деплоя в Kubernetes
	fmt.Printf("Deploying %s\n", name)
	// Здесь должна быть логика создания ресурсов Kubernetes
	return nil
}

func delete(name string) error {
	// Логика удаления из Kubernetes
	fmt.Printf("Deleting %s\n", name)
	ctx := context.Background()
	err := clientset.CoreV1().Pods("axis-testing").Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
