package main

import (
	"net/http"

	"github.com/bitsbeats/drone-tree-config/plugin"

	"github.com/drone/drone-go/plugin/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type (
	spec struct {
		Concat              bool   `envconfig:"PLUGIN_CONCAT"`
		MaxDepth            int    `envconfig:"PLUGIN_MAXDEPTH" default:"2"`
		Fallback            bool   `envconfig:"PLUGIN_FALLBACK"`
		Debug               bool   `envconfig:"PLUGIN_DEBUG"`
		WhitelistFile       string `envconfig:"PLUGIN_WHITELIST_FILE"`
		Address             string `envconfig:"PLUGIN_ADDRESS" default:":3000"`
		Secret              string `envconfig:"PLUGIN_SECRET"`
		GitHubToken         string `envconfig:"GITHUB_TOKEN"`
		Server              string `envconfig:"SERVER" default:"https://api.github.com"`
		BitBucketAuthServer string `envconfig:"BITBUCKET_AUTH_SERVER"`
		BitBucketClient     string `envconfig:"BITBUCKET_CLIENT"`
		BitBucketSecret     string `envconfig:"BITBUCKET_SECRET"`
	}
)

func main() {
	spec := new(spec)
	if err := envconfig.Process("", spec); err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.GitHubToken == "" && (spec.BitBucketClient == "" || spec.BitBucketSecret == "") {
		logrus.Warnln("missing SCM credentials, e.g. GitHub token")
	}
	if spec.Address == "" {
		spec.Address = ":3000"
	}
	if spec.BitBucketAuthServer == "" {
		spec.BitBucketAuthServer = spec.Server
	}

	handler := config.Handler(
		plugin.New(
			plugin.WithBitBucketAuthServer(spec.BitBucketAuthServer),
			plugin.WithBitBucketClient(spec.BitBucketClient),
			plugin.WithBitBucketSecret(spec.BitBucketSecret),
			plugin.WithConcat(spec.Concat),
			plugin.WithFallback(spec.Fallback),
			plugin.WithGithubToken(spec.GitHubToken),
			plugin.WithMaxDepth(spec.MaxDepth),
			plugin.WithServer(spec.Server),
			plugin.WithWhitelistFile(spec.WhitelistFile),
		),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Address)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Address, nil))
}
