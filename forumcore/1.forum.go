package forumcore

import (
	"embed"
	"errors"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/FimGroup/fim/components"
	"github.com/FimGroup/fim/fimapi/basicapi"
	"github.com/FimGroup/fim/fimcore"
	"github.com/FimGroup/fim/fimsupport/logging"
)

//go:embed flowmodel.*.toml
var flowModelFs embed.FS

//go:embed scene.*.toml
var sceneFs embed.FS

func StartForum() error {
	// init logging
	{
		lm, err := logging.NewLoggerManager("logs/forum", 7, 20*1024*1024, 5, logrus.InfoLevel, true, false)
		if err != nil {
			return err
		}
		logging.SetLoggerManager(lm)
	}
	// init fim core
	if err := fimcore.Init(); err != nil {
		return err
	}
	// create container
	container := fimcore.NewUseContainer()
	// init plugins/components
	if err := components.InitComponent(container); err != nil {
		return err
	}
	// setup configure manager
	settableConfigureManager := fimcore.NewSettableConfigureManager()
	{
		dburl, ok := os.LookupEnv("DATABASE_URL")
		if !ok {
			panic(errors.New("database url is not set in env"))
		}
		settableConfigureManager.SetConfigure("forum_database", dburl)
	}
	// load configure manager
	if err := loadConfigureManager(container, []basicapi.ConfigureManager{
		fimcore.NewEnvConfigureManager(),
		settableConfigureManager,
	}); err != nil {
		return err
	}
	// load custom functions
	if err := loadCustomFn(container, map[string]basicapi.FnGen{
		"#print_obj": FnPrintObject,
		"#panic":     FnPanic,
	}); err != nil {
		return err
	}

	// load FlowModels
	if err := loadFlowModel(container, []string{
		"flowmodel.all.toml",
	}); err != nil {
		return err
	}
	// load pipelines/flows
	if err := loadMerged(container, []string{
		"scene.user.register.toml",
		"scene.user.login.toml",
		"scene.forums.new_forum.toml",
		"scene.posts.new_post.toml",
		"scene.posts.list_posts_by_forum.toml",
	}); err != nil {
		return err
	}

	// start container
	if err := container.StartContainer(); err != nil {
		return err
	}

	return nil
}

func loadConfigureManager(container basicapi.BasicContainer, managers []basicapi.ConfigureManager) error {
	for _, v := range managers {
		if err := container.AddConfigureManager(v); err != nil {
			return err
		}
	}
	return nil
}

func loadCustomFn(container basicapi.BasicContainer, mapping map[string]basicapi.FnGen) error {
	for name, fg := range mapping {
		if err := container.RegisterCustomFn(name, fg); err != nil {
			return err
		}
	}
	return nil
}

func loadFlowModel(container basicapi.BasicContainer, files []string) error {
	for _, file := range files {
		data, err := flowModelFs.ReadFile(file)
		if err != nil {
			return err
		}
		log.Println("read FlowModel content:", string(data))

		if err := container.LoadFlowModel(string(data)); err != nil {
			return err
		}
	}
	return nil
}

func loadMerged(container basicapi.BasicContainer, files []string) error {
	for _, file := range files {
		data, err := sceneFs.ReadFile(file)
		if err != nil {
			return err
		}
		log.Println("read scene content:", string(data))

		if err := container.LoadMerged(string(data)); err != nil {
			return err
		}
	}
	return nil
}
