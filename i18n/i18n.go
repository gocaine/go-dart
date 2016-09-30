package i18n

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/olebedev/config"
	"golang.org/x/text/language"
)

var supportedLanguages = []language.Tag{
	language.AmericanEnglish,
	language.French,
}

var matcher = language.NewMatcher(supportedLanguages)

var lclzr *localizer

type localizer struct {
	trans map[string]*config.Config
	base  *config.Config
}

func init() {
	lclzr = new(localizer)
	lclzr.trans = make(map[string]*config.Config)

	// load base
	lclzr.base = baseConfig
}

// Translations serves the requested translation if it exists
func Translations(locale string) (json string, err error) {
	if tag, _, confidence := matcher.Match(language.Make(locale)); confidence >= language.High {
		baseTag, _ := tag.Base()
		shortLocale := baseTag.ISO3()

		var cfg *config.Config

		if shortLocale == "eng" {
			cfg = lclzr.base
		} else {
			if alreadyLoadedConf, ok := lclzr.trans[shortLocale]; ok {
				cfg = alreadyLoadedConf
			} else {
				var tmpConf *config.Config
				tmpConf, err = config.ParseYamlFile("messages." + shortLocale + ".yaml")
				if err != nil {
					return
				}
				cfg, err = lclzr.base.Extend(tmpConf)
				if err != nil {
					return
				}
				lclzr.trans[shortLocale] = cfg
			}
		}

		json, err = config.RenderJson(cfg)
	} else {
		err = errors.New("Not supported or recognized language")
	}
	return
}

// BaseTranslation gives the default translation for the given key
func BaseTranslation(key string) (string, bool) {
	str, err := lclzr.base.String(key)
	if err != nil {
		log.Error(err)
		return "", false
	}
	return str, true
}
