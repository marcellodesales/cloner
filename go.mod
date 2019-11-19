module github.com/marcellodesales/cloner

go 1.13

require (
	github.com/jinzhu/copier v0.0.0-20190625015134-976e0346caa8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/prometheus/common v0.4.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/thoas/go-funk v0.4.0
	github.com/uber/prototool v1.8.1
	gopkg.in/yaml.v2 v2.2.2
)

replace google.golang.org/genproto v0.0.0-20170818100345-ee236bd376b0 => google.golang.org/genproto v0.0.0-20170818010345-ee236bd376b0
