package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var ServerConf *Conf

type Conf struct {
	GameIp  string `yaml:"GameIp"`
	SockBuf int    `yaml:"SockBuf"`
}

type SeverInfo struct {
	NodeId      int    `yaml:"NodeId"`
	Name        string `yaml:"Name"`
	IpPort      string `yaml:"IpPort"`
	EtcdVersion string `yaml:"EtcdVersion"`
}

func (c *Conf) getServerConf() *Conf {
	path, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(path + `\server_conf.yaml`)
	if os.IsNotExist(err) {
		_, err = os.Create("server_conf.yaml")
		if err != nil {
			panic("os.Create(\"server_conf.yaml\") err:" + err.Error())
		}
		yamlFile, err = ioutil.ReadFile("server_conf.yaml")
	}
	if err != nil {
		panic("os.ReadFile(\"server_conf.yaml\") err:" + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic("yaml.Unmarshal(yamlFile, c) err:" + err.Error())
	}
	return c
}

func (c *Conf) SaveServerConf() {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		panic("yaml.Marshal(c) err:" + err.Error())
	}
	err = ioutil.WriteFile("server_conf.yaml", bytes, os.ModePerm)
	if err != nil {
		panic("ioutil.WriteFile err:" + err.Error())
	}
}

func init() {
	ServerConf = (&Conf{}).getServerConf()
	ServerConf.GameIp = `:9999`
	ServerConf.SockBuf = 1024 * 1024 //1m 缓存区大小
	ServerConf.SaveServerConf()
}
