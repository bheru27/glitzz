package modules

import (
	"fmt"
	"github.com/bheru27/glitzz/config"
	"github.com/bheru27/glitzz/core"
	"github.com/bheru27/glitzz/logging"
	"github.com/bheru27/glitzz/modules/c3"
	"github.com/bheru27/glitzz/modules/decide"
	"github.com/bheru27/glitzz/modules/degeneracy"
	"github.com/bheru27/glitzz/modules/fourchan"
	"github.com/bheru27/glitzz/modules/info"
	"github.com/bheru27/glitzz/modules/links"
	"github.com/bheru27/glitzz/modules/pipes"
	"github.com/bheru27/glitzz/modules/quotes"
	"github.com/bheru27/glitzz/modules/reactions"
	"github.com/bheru27/glitzz/modules/reddit"
	"github.com/bheru27/glitzz/modules/reminders"
	"github.com/bheru27/glitzz/modules/sed"
	"github.com/bheru27/glitzz/modules/seen"
	"github.com/bheru27/glitzz/modules/stackexchange"
	"github.com/bheru27/glitzz/modules/tell"
	"github.com/bheru27/glitzz/modules/tv"
	"github.com/bheru27/glitzz/modules/untappd"
	"github.com/bheru27/glitzz/modules/vatsim"
	"github.com/pkg/errors"
)

var log = logging.New("modules")

type moduleConstructor func(sender core.Sender, conf config.Config) (core.Module, error)

func CreateModules(sender core.Sender, conf config.Config) ([]core.Module, error) {
	return createModules(getModuleConstructors(), conf.EnabledModules, sender, conf)
}

func createModules(moduleConstructors map[string]moduleConstructor, moduleNames []string, sender core.Sender, conf config.Config) ([]core.Module, error) {
	var modules []core.Module
	for _, moduleName := range conf.EnabledModules {
		moduleConstructor, ok := moduleConstructors[moduleName]
		if !ok {
			return nil, fmt.Errorf("module %s does not exist and could not be loaded", moduleName)
		}
		module, err := moduleConstructor(sender, conf)
		if err != nil {
			return nil, errors.Wrapf(err, "module %s could not be created", moduleName)
		}
		log.Info("created module", "name", moduleName)
		modules = append(modules, module)

	}
	return modules, nil
}

func getModuleConstructors() map[string]moduleConstructor {
	modules := map[string]moduleConstructor{
		"c3":            c3.New,
		"decide":        decide.New,
		"degeneracy":    degeneracy.New,
		"fourchan":      fourchan.New,
		"info":          info.New,
		"links":         links.New,
		"pipes":         pipes.New,
		"quotes":        quotes.New,
		"reactions":     reactions.New,
		"reddit":        reddit.New,
		"reminders":     reminders.New,
		"stackexchange": stackexchange.New,
		"sed":           sed.New,
		"seen":          seen.New,
		"tell":          tell.New,
		"tv":            tv.New,
		"untappd":       untappd.New,
		"vatsim":        vatsim.New,
	}
	return modules
}
