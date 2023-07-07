package forumcore

import (
	"embed"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/FimGroup/fim/components"
	"github.com/FimGroup/fim/fimapi/basicapi"
	"github.com/FimGroup/fim/fimapi/pluginapi"
	"github.com/FimGroup/fim/fimcore"
	"github.com/FimGroup/fim/fimsupport/resourcemanager"
	"github.com/FimGroup/logging"

	"github.com/sirupsen/logrus"
)

//go:embed flowmodel.*.toml
var flowModelFs embed.FS

//go:embed scene.*.toml
var sceneFs embed.FS

//go:embed connector.*.toml
var connectorFs embed.FS

func StartForum() error {
	// init logging
	{
		lm, err := logging.NewLoggerManager("logs/forum", 7, 20*1024*1024, 5, logrus.InfoLevel, true, false)
		if err != nil {
			return err
		}
		logging.SetLoggerManager(lm)
	}
	// init file manager
	templateFileManager := resourcemanager.NewOsFileResourceManager("template_file_manager", filepath.Join("web", "templates"))
	// setup configure manager
	settableConfigureManager := fimcore.NewSettableConfigureManager()
	{
		dburl, ok := os.LookupEnv("DATABASE_URL")
		if !ok {
			panic(errors.New("database url is not set in env"))
		}
		settableConfigureManager.SetConfigure("forum_database", dburl)
	}

	// init fim package and application
	if err := fimcore.Init(); err != nil {
		return err
	}
	app := fimcore.NewPluginApplication()
	// setup app configure manager
	if err := loadAppConfigureManager(app, []basicapi.FullConfigureManager{
		fimcore.NewEnvConfigureManager(),
		settableConfigureManager,
	}); err != nil {
		return err
	}
	// setup file manager
	if err := app.AddFileResourceManager(templateFileManager); err != nil {
		return err
	}
	// setup sub connectors
	if err := loadAppSubConnectors(app, []string{
		"connector.shared.toml",
	}); err != nil {
		return err
	}
	// load connectors
	if err := components.InitConnectors(app); err != nil {
		return err
	}
	if err := app.Startup(); err != nil {
		return err
	}

	// create container
	container := app.SpawnUseContainer()
	// init plugins/components
	if err := components.InitFunctions(container); err != nil {
		return err
	}
	// load configure manager
	if err := loadConfigureManager(container, []basicapi.ConfigureManager{
		fimcore.NewEnvConfigureManager(),
		settableConfigureManager,
	}); err != nil {
		return err
	}
	// load custom functions
	customFunctions := &CustomFunctions{_logger: logging.GetLoggerManager().GetLogger("Forum")}
	if err := loadCustomFn(container, map[string]basicapi.FnGen{
		"#print_obj": customFunctions.FnPrintObject,
		"#panic":     customFunctions.FnPanic,
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
		"scene.show_user_register.toml",
	}); err != nil {
		return err
	}

	// start container
	if err := container.StartContainer(); err != nil {
		return err
	}

	return nil
}

func loadAppConfigureManager(app pluginapi.ApplicationSupport, managers []basicapi.FullConfigureManager) error {
	for _, v := range managers {
		if err := app.AddConfigureManager(v); err != nil {
			return err
		}
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

func loadAppSubConnectors(app pluginapi.ApplicationSupport, files []string) error {
	for _, file := range files {
		data, err := connectorFs.ReadFile(file)
		if err != nil {
			return err
		}
		log.Println("read app sub content:", string(data))

		if err := app.AddSubConnectorGeneratorDefinitions(string(data)); err != nil {
			return err
		}
	}
	return nil
}
