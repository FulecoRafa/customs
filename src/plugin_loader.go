package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"plugin"

	"github.com/FulecoRafa/customs/lib"
)

type PluginSetup struct {
    Name string `json:"name"`
    Path string `json:"path"`
}

type Config struct {
	Plugins []PluginSetup  `json:"plugins"`
}

var flagPlugins []func()

var prePlugins []lib.ReqHookFunc
var postPlugins []lib.ResHookFunc

func loadConfigFile(filePath string) (Config, error) {
    errStr := "Error reading config file: %w"
    file, err := os.ReadFile(filePath)
    if err != nil {
        return Config{}, fmt.Errorf(errStr, err)
    }
    var result Config
    err = json.Unmarshal(file, &result)
    if err != nil {
        return Config{}, fmt.Errorf(errStr, err)
    }
    return result, nil
}

func LoadPlugins(config_path string) error {
    config, err := loadConfigFile(config_path)
    if err != nil {
        return err
    }

    pluginN:= len(config.Plugins)
    flagPlugins = make([]func(), 0, pluginN)
    prePlugins = make([]lib.ReqHookFunc, 0, pluginN)
    postPlugins = make([]lib.ResHookFunc, 0, pluginN)

    for _, p := range config.Plugins {
        loaded, err := loadPlugin(p)
        if err != nil {
            slog.Warn("Failed to load plugin, skipping...", "plugin", p.Name, "error", err)
        }
        flagPlugins = append(flagPlugins, loaded.AddConfigFlag)
        prePlugins = append(prePlugins, loaded.PreRequestHook)
        postPlugins = append(postPlugins, loaded.PostRequestHook)
    }

    return nil
}

const pluginSymbol string = "Plugin"
func loadPlugin(p PluginSetup) (lib.CustomsPlugin, error) {
    pf, err := plugin.Open(p.Path)
    if err != nil {
        return nil, err
    }
    symb, err := pf.Lookup(pluginSymbol)
    cp, ok := symb.(lib.CustomsPlugin)
    if !ok {
        return nil, fmt.Errorf("Declared plugin does not export a `Plugin` variable of type `CustomsPlugin`")
    }
    return cp, nil
}

func RegisterPluginFlags() {
    for _, f := range flagPlugins {
        f()
    }
}

func RunPreRequestHooks(req *http.Request, r lib.Redirect) {
    for _, f := range prePlugins {
        f(req, r)
    }
}

func RunPostRequestHooks(res *http.Response, r lib.Redirect) {
    for _, f := range postPlugins {
        f(res, r)
    }
}
