package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"k3s-client/app/models"
	"k3s-client/global"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

type containerServices struct {
}

var ContainerServices = new(containerServices)

type ChallengeRequest struct {
	Username    string `json:"username"`
	ChallengeID uint   `json:"challenge_id"`
}

type UserContainerResponse struct {
	ChallengeName string `json:"name"`         // 挑战名称
	ContainerID   uint   `json:"container_id"` // 容器名称
	Namespace     string `json:"namespace"`    // 命名空间
	Flag          string `json:"flag"`         // flag
	Status        string `json:"status"`       // 容器状态
	Url           string `json:"url"`          // 容器访问地址
}

func (containerServices *containerServices) GetUserContainersInfo(id uint) (err error, response UserContainerResponse) {
	// intId, err := strconv.Atoi(id)
	var containers models.UserContainer
	if err := global.App.DB.First(&containers, id).Error; err != nil {
		return err, response
	}

	var template models.ChallengeTemplate
	if err := global.App.DB.First(&template, containers.TemplateID).Error; err != nil {
		return err, response
	}

	response = UserContainerResponse{
		ChallengeName: template.Name,
		ContainerID:   containers.ID.ID,
		Namespace:     containers.Namespace,
		Flag:          containers.Flag,
		Status:        containers.Status,
		Url:           containers.Url,
	}

	return
}

type CreatedResource struct {
	ResourceInterface dynamic.ResourceInterface
	Name              string
}

func (containerServices *containerServices) CreateContainer(templateID uint, c *gin.Context) (err error, response UserContainerResponse) {
	var template models.ChallengeTemplate
	if err := global.App.DB.First(&template, templateID).Error; err != nil {
		// err = errors.New("Challenge template not found")
		return err, response
	}
	userId := c.Keys["userId"].(uint)

	var containers models.UserContainer
	if err := global.App.DB.Where("user_id = ?", userId).First(&containers).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err, response
		} else {

		}
	} else {
		if containers.TemplateID == templateID {
			err = fmt.Errorf("%d'%s already exists", userId, template.Name)
			return err, response
		}
	}
	namespace := fmt.Sprintf("ctf-%d", userId)

	_, err = global.App.Client.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		// NameSpace not exists
		// NameSpace Create
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = global.App.Client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			return err, response
		}
		// return err, response
	}

	flag := "CNSS{" + uuid.New().String() + "}"

	hasher := sha256.New()

	// 将时间戳和 userId 写入哈希对象
	hasher.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatUint(uint64(userId), 10)))

	// 获取哈希值
	prefixHash := hex.EncodeToString(hasher.Sum(nil))
	prefix := prefixHash[:32] + "." + template.Name
	name := template.Name + "-" + time.Now().Format("20060102150405")
	// 替换占位符
	templateString := string(template.TemplateYAML)
	replacements := map[string]string{
		"<NAME>":        name,
		"<FLAG>":        flag,
		"<PREFIX-HASH>": prefix,
	}

	for placeholder, value := range replacements {
		templateString = strings.ReplaceAll(templateString, placeholder, value)
	}

	// var deployment v1.Deployment
	// if err := yamlutil.Unmarshal([]byte(template.ConfigYAML), &deployment); err != nil {
	// 	err = errors.New("Failed to unmarshal deployment YAML")
	// 	return err, ""
	// }

	// deployment.Namespace = namespace
	// _, err = global.App.Client.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	// if err == nil {
	// 	return err, ""
	// }
	// log.Default().Println(templateString)
	createdResources := []CreatedResource{}
	var dri dynamic.ResourceInterface
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(templateString)), 100)
	for {
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			if err.Error() == "EOF" { // End of file reached, which means decoding is complete
				break
			}
			return
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return err, response
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			global.App.Log.Error("error", zap.Any("error", err))
			return err, response
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(global.App.Client.Discovery())
		if err != nil {
			global.App.Log.Error("error", zap.Any("error", err))
			return err, response
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			global.App.Log.Error("error", zap.Any("error", err))
			return err, response
		}

		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespace)
			}
			dri = global.App.DynamicClient.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = global.App.DynamicClient.Resource(mapping.Resource)
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			global.App.Log.Error("error", zap.Any("error", err))
			return err, response
		}
		createdResources = append(createdResources, CreatedResource{
			ResourceInterface: dri,
			Name:              obj2.GetName(),
		})
		// log.Println(resourceCreated)
		// global.App.Log.Info("created:", zap.Any("kind", obj2.GetKind()), zap.Any("name", obj2.GetName()))
	}
	log.Println(createdResources)
	// pod, err := global.App.DynamicClient.Resource(corev1.SchemeGroupVersion.WithResource("pods")).Namespace(namespace).Get(context.TODO(), resourceCreated[2], v1.GetOptions{})
	// //global.App.Client.CoreV1().Pods(namespace).Get(context.TODO(), obj2Name, v1.GetOptions{})
	// if err != nil {
	// 	deleteErr := deleteResources( resourceCreated)
	// 	if deleteErr != nil {
	// 		global.App.Log.Error("Failed to delete Kubernetes resource after database insertion failure", zap.Any("error", deleteErr))
	// 	}
	// 	return err, response
	// }
	// log.Default().Println(pod.Object["status"].(map[string]interface{})["phase"].(string))

	pods, err := global.App.Client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "ctf=" + name,
	})
	if err != nil {
		deleteErr := deleteResources(createdResources)
		if deleteErr != nil {
			global.App.Log.Error("Failed to delete Kubernetes resource after database insertion failure", zap.Any("error", deleteErr))
			return deleteErr, response
		}
		return err, response
	}

	userContainer := models.UserContainer{
		Name:       template.Name,
		Namespace:  namespace,
		Flag:       flag,
		Status:     string(pods.Items[0].Status.Phase),
		Url:        prefix + ".ctf.fanllspd.com",
		UserID:     userId,
		TemplateID: templateID,
		Timestamps: models.Timestamps{
			UpdateAt:  time.Now(),
			CreatedAt: time.Now(),
		},
	}
	if err := global.App.DB.Create(&userContainer).Error; err == nil {
		deleteErr := deleteResources(createdResources)
		if deleteErr != nil {
			global.App.Log.Error("Failed to delete Kubernetes resource after database insertion failure", zap.Any("error", deleteErr))
			return deleteErr, response
		}
		return err, response
	}

	// 返回容器信息
	return nil, UserContainerResponse{
		ChallengeName: template.Name,
		ContainerID:   userContainer.ID.ID,
		Namespace:     userContainer.Namespace,
		Flag:          userContainer.Flag,
		Status:        userContainer.Status,
		Url:           userContainer.Url,
	}
}

func deleteResources(resources []CreatedResource) error {
	for _, resource := range resources {
		err := resource.ResourceInterface.Delete(context.Background(), resource.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
