package forumcore

import (
	"embed"
	"log"

	"github.com/FimGroup/fim/components"
	"github.com/FimGroup/fim/fimapi/basicapi"
	"github.com/FimGroup/fim/fimcore"
)

//go:embed flowmodel.*.toml
var flowModelFs embed.FS

//go:embed scene.*.toml
var sceneFs embed.FS

func StartForum() error {
	container := fimcore.NewUseContainer()
	if err := components.InitComponent(container); err != nil {
		return err
	}
	if err := loadCustomFn(container, map[string]basicapi.FnGen{
		"#print_obj": FnPrintObject,
		"#panic":     FnPanic,
	}); err != nil {
		return err
	}

	if err := loadFlowModel(container, []string{
		"flowmodel.all.toml",
	}); err != nil {
		return err
	}
	if err := loadMerged(container, []string{
		"scene.user.register.toml",
		"scene.user.login.toml",
		"scene.forums.new_forum.toml",
		"scene.posts.new_post.toml",
	}); err != nil {
		return err
	}

	if err := container.StartContainer(); err != nil {
		return err
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
