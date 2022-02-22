package configs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/kubefilebrowser/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"net"
	"os"
	"path/filepath"
)

// Configure stores configuration.
type Configure struct {
	RunMode     string   `envconfig:"RUN_MODE" default:"debug"`
	HTTPPort    string   `envconfig:"HTTP_PORT" default:"9999"`
	HTTPAddr    string   `envconfig:"HTTP_ADDR" default:"0.0.0.0"`
	IPWhiteList []string `envconfig:"IP_WHITE_LIST" default:"*"`
	RootPath    string   `envconfig:"ROOT_PATH" default:""`
}

var (
	TmpPath       = os.TempDir()
	Config        Configure
	RestClient    *kubernetes.Clientset
	KuBeResConf   *rest.Config
	envFile       = kingpin.Flag("env_file", "Load the environment variable file").Default(".envfile").String()
	rootPath      = kingpin.Flag("root_path", "Save data directory").Default("").String()
)

const notFoundKubeConfig = `Missing or incomplete kubernetes configuration info.  Please point to an existing, complete config file:

  1. Via the KUBECONFIG environment variable
  2. In your home directory as ~/.kube/config`

func LoadConfig() {
	logrus.Debug("Load variable")
	// load environment variables from file.
	_ = godotenv.Load(*envFile)

	// load the configuration from the environment.
	err := envconfig.Process("", &Config)
	if err != nil {
		logrus.Fatal(err)
	}
	if *rootPath != "" {
		Config.RootPath = *rootPath
	}
	if !utils.InSliceString("*", Config.IPWhiteList) {
		for _, ip := range Config.IPWhiteList {
			if net.ParseIP(ip) != nil {
				continue
			}
			logrus.Fatal(fmt.Sprint(ip, ", Invalid whitelist"))
		}
	}
	KuBeResConf, err = kConfig()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
	RestClient, err = InitRestClient()
	if err != nil {
		fmt.Println(notFoundKubeConfig)
		os.Exit(1)
	}
}

func kConfig() (conf *rest.Config, err error) {
	if Config.RunMode == gin.DebugMode {
		home, _ := homedir.Dir()
		var kubeConfigEnv = os.Getenv("KUBECONFIG")
		var kubeConfig = filepath.Join(home, ".kube", "config")
		var kuBeConf []byte
		if utils.FileOrPathExist(kubeConfigEnv) {
			kuBeConf, err = ioutil.ReadFile(kubeConfigEnv)
		} else if utils.FileOrPathExist(kubeConfig) {
			kuBeConf, err = ioutil.ReadFile(kubeConfig)
		}
		if err != nil {
			return nil, err
		}
		conf, err = clientcmd.RESTConfigFromKubeConfig(kuBeConf)
	} else {
		conf, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	conf.Timeout = 0
	return
}

func InitRestClient() (*kubernetes.Clientset, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	kConf.QPS = 500
	kConf.Burst = 1000
	return kubernetes.NewForConfig(kConf)
}

func InitMetricsClient() (metrics.Interface, error) {
	kConf, err := kConfig()
	if err != nil {
		return nil, err
	}
	return metrics.NewForConfig(kConf)
}
