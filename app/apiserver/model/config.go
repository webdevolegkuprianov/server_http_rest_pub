package model

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//config yaml struct
type Service struct {
	APIVersion string `yaml:"apiVersion"`
	Spec       struct {
		Ports struct {
			Name string `yaml:"name"`
			Addr string `yaml:"bind_addr"`
		} `yaml:"ports"`
		DBpg struct {
			Name              string `yaml:"name"`
			Host              string `yaml:"host"`
			Port              uint16 `yaml:"port"`
			User              string `yaml:"user"`
			Password          string `yaml:"password"`
			Database          string `yaml:"database"`
			MaxConnLifetime   int    `yaml:"max_conn_lifetime"`
			MaxConnIdletime   int    `yaml:"max_conn_idletime"`
			MaxConns          int32  `yaml:"max_conns"`
			MinConns          int32  `yaml:"min_conns"`
			HealthCheckPeriod int    `yaml:"health_check_period"`
		} `yaml:"dbpg"`
		DBms struct {
			Name string `yaml:"name"`
			Url  string `yaml:"url"`
		} `yaml:"dbms"`
		Jwt struct {
			TokenDecode string `yaml:"token"`
			LifeTerm    int    `yaml:"term"`
		} `yaml:"jwt"`
		Client struct {
			UrlGazCrmTest     string `yaml:"url_gaz_crm_test"`
			UrlMailingService string `yaml:"url_mailing_service"`
		} `yaml:"client"`
		Queryies struct {
			Booking          string `yaml:"booking"`
			Stocks           string `yaml:"stocks"`
			BasicModelsPrice string `yaml:"basic_models_price"`
			OptionsPrice     string `yaml:"options_price"`
			GeneralPrice     string `yaml:"general_price"`
			Sprav            string `yaml:"sprav"`
			Sprav_new        string `yaml:"sprav_new"`
			Options          string `yaml:"options"`
			OptionsSprav     string `yaml:"options_sprav"`
			Packets          string `yaml:"packets"`
			Colors           string `yaml:"colors"`
		} `yaml:"queryies"`
	} `yaml:"spec"`
}

//New config
func NewConfig() (*Service, error) {

	var service *Service

	f, err := filepath.Abs("/root/config/server_http_rest.yaml")
	if err != nil {
		return nil, err
	}

	y, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(y, &service); err != nil {
		return nil, err
	}

	return service, nil

}
