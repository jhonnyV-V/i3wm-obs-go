package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/goccy/go-yaml"
	"go.i3wm.org/i3"
)

type Config struct {
	Password           string   `yaml:"password"`
	Url                string   `yaml:"url"`
	SourceName         string   `yaml:"sourceName"`
	ExcludedWorkspaces []string `yaml:"excludedWorkspaces"`
	IsScene            bool     `yaml:"isScene"`
}

func getItemsFromScene(name string, client *goobs.Client) (*sceneitems.GetSceneItemListResponse, error) {
	listParams := sceneitems.NewGetSceneItemListParams().WithSceneName(name)
	items, err := client.SceneItems.GetSceneItemList(listParams)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func setItemEnabled(
	enabled bool,
	sceneName string,
	item *typedefs.SceneItem,
	client *goobs.Client,
) error {
	sceneParams := sceneitems.NewSetSceneItemEnabledParams().WithSceneItemEnabled(enabled).WithSceneItemId(item.SceneItemID).WithSceneName(sceneName)
	_, err := client.SceneItems.SetSceneItemEnabled(sceneParams)
	return err
}

func setAll(
	enabled bool,
	sceneName string,
	items *sceneitems.GetSceneItemListResponse,
	client *goobs.Client,
) {
	for _, v := range items.SceneItems {
		err := setItemEnabled(enabled, sceneName, v, client)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//ignore case
			return false
		} else {
			log.Printf("failed to get file %s: %s\n", path, err)
			return false
		}
	}
	return true
}

func getConfigPath() string {
	userConfigPath, err := os.UserConfigDir()
	if err != nil {
		log.Printf("failed to get config dir: %s\n", err)
		os.Exit(1)
	}

	return userConfigPath + "/i3wn-obs"
}

func CreateConfigDir() string {
	configPath := getConfigPath()
	_, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(configPath, 0766)
			if err != nil {
				log.Printf("failed to create i3wn-obs dir: %s\n", err)
				os.Exit(1)
			}
		} else {
			log.Printf("failed to get i3wn-obs dir: %s\n", err)
			os.Exit(1)
		}
	}

	return configPath
}

func createConfig(configPath string) error {
	configYaml := fmt.Sprintf(
		"%s/config.yaml",
		configPath,
	)

	defaultConfig := Config{
		Url:        "localhost:4455",
		Password:   "SecurePassword",
		SourceName: "ThingToBlock",
		IsScene:    false,
	}

	parsed, err := yaml.Marshal(defaultConfig)
	if err != nil {
		log.Printf("Failed to parse default config: %s\n", err)
		return err
	}

	err = os.WriteFile(configYaml, parsed, 0666)
	if err != nil {
		log.Printf("Failed to create config: %s\n", err)
		return err
	}
	return nil
}

func readConfig(configPath string) Config {
	configYaml := fmt.Sprintf(
		"%s/config.yaml",
		configPath,
	)

	if !fileExist(configYaml) {
		if err := createConfig(configPath); err != nil {
			os.Exit(1)
		}
		log.Printf("Edit the config file %s to add your password and to be able to fit your use\n", configYaml)
		os.Exit(0)
	}

	file, err := os.ReadFile(configYaml)
	if err != nil {
		log.Printf("failed to read from config %v\n", err)
		os.Exit(1)
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Printf("failed to parse from config %v\n", err)
		os.Exit(1)
	}

	return config
}

func main() {
	recv := i3.Subscribe(i3.WorkspaceEventType)
	userConf := readConfig(CreateConfigDir())

	client, err := goobs.New(userConf.Url, goobs.WithPassword(userConf.Password))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	sceneParams := &scenes.GetCurrentProgramSceneParams{}
	currentScene, err := client.Scenes.GetCurrentProgramScene(sceneParams)
	if err != nil {
		log.Fatal(err)
	}
	sceneName := currentScene.SceneName

	for recv.Next() {
		ev := recv.Event().(*i3.WorkspaceEvent)
		if ev.Change != "focus" {
			continue
		}
		if slices.Contains(userConf.ExcludedWorkspaces, ev.Current.Name) {
			items, err := getItemsFromScene(sceneName, client)
			if err != nil {
				log.Fatal(err)
			}

			if userConf.IsScene || userConf.SourceName == "" {
				setAll(false, sceneName, items, client)
				continue
			}

			if userConf.SourceName != "" {
				for _, v := range items.SceneItems {
					if v.SourceName == userConf.SourceName {
						err = setItemEnabled(false, sceneName, v, client)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				continue
			}
		}

		if slices.Contains(userConf.ExcludedWorkspaces, ev.Old.Name) && !slices.Contains(userConf.ExcludedWorkspaces, ev.Current.Name) {
			items, err := getItemsFromScene(sceneName, client)
			if err != nil {
				log.Fatal(err)
			}

			if userConf.IsScene || userConf.SourceName == "" {
				setAll(true, sceneName, items, client)
				continue
			}

			if userConf.SourceName != "" {
				for _, v := range items.SceneItems {
					if v.SourceName == userConf.SourceName {
						err = setItemEnabled(true, sceneName, v, client)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				continue
			}
		}
	}
	log.Fatal(recv.Close())
}
