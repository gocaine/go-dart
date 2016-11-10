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

// GetLocale normalize locale to ISO3
func GetLocale(langHeader string) string {
	log.WithField("header", langHeader).Info("GetLocale")
	baseTag, _ := language.English.Base()
	tags, _, err := language.ParseAcceptLanguage(langHeader)
	if err != nil {
		log.Error(err)
		return baseTag.ISO3()
	}
	if tag, _, confidence := matcher.Match(tags...); confidence >= language.High {
		baseTag, _ = tag.Base()
	}
	log.WithField("baseTag", baseTag).Debug("GetLocale")
	return baseTag.ISO3()
}

func _translations(locale string) (cfg *config.Config, err error) {
	if tag, _, confidence := matcher.Match(language.Make(locale)); confidence >= language.High {
		baseTag, _ := tag.Base()
		shortLocale := baseTag.ISO3()
		cfg = lclzr.base

		if shortLocale != "eng" {
			if alreadyLoadedConf, ok := lclzr.trans[shortLocale]; ok {
				cfg = alreadyLoadedConf
			} else {
				var tmpConf *config.Config
				yamlData, ok := yamlFiles[shortLocale]
				if !ok {
					return
				}
				tmpConf, err = config.ParseYaml(yamlData)
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
	} else {
		err = errors.New("Not supported or recognized language")
	}
	return
}

// Translation serves the translation for the given key and locale
func Translation(key, locale string) string {
	log.WithField("key", key).WithField("locale", locale).Info("Translation")
	cfg, err := _translations(locale)
	if err != nil {
		log.Error("Error in Translation/_translation", err)
		return ""
	}
	str, err := cfg.String(key)
	if err != nil {
		log.Error("Error in Translation/missing key ?", key, err)
		return ""
	}
	return str
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
