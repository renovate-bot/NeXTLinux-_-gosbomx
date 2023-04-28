package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/scylladb/go-set/strset"

	"github.com/nextlinux/syft/syft/file"
)

const ApkDBGlob = "**/lib/apk/db/installed"

var _ FileOwner = (*ApkMetadata)(nil)

// ApkMetadata represents all captured data for a Alpine DB package entry.
// See the following sources for more information:
// - https://wiki.alpinelinux.org/wiki/Apk_spec
// - https://git.alpinelinux.org/apk-tools/tree/src/package.c
// - https://git.alpinelinux.org/apk-tools/tree/src/database.c
type ApkMetadata struct {
	Package       string          `mapstructure:"P" json:"package"`
	OriginPackage string          `mapstructure:"o" json:"originPackage" cyclonedx:"originPackage"`
	Maintainer    string          `mapstructure:"m" json:"maintainer"`
	Version       string          `mapstructure:"V" json:"version"`
	License       string          `mapstructure:"L" json:"license"`
	Architecture  string          `mapstructure:"A" json:"architecture"`
	URL           string          `mapstructure:"U" json:"url"`
	Description   string          `mapstructure:"T" json:"description"`
	Size          int             `mapstructure:"S" json:"size" cyclonedx:"size"`
	InstalledSize int             `mapstructure:"I" json:"installedSize" cyclonedx:"installedSize"`
	Dependencies  []string        `mapstructure:"D" json:"pullDependencies" cyclonedx:"pullDependencies"`
	Provides      []string        `mapstructure:"p" json:"provides" cyclonedx:"provides"`
	Checksum      string          `mapstructure:"C" json:"pullChecksum" cyclonedx:"pullChecksum"`
	GitCommit     string          `mapstructure:"c" json:"gitCommitOfApkPort" cyclonedx:"gitCommitOfApkPort"`
	Files         []ApkFileRecord `json:"files"`
}

type spaceDelimitedStringSlice []string

func (m *ApkMetadata) UnmarshalJSON(data []byte) error {
	var fields []reflect.StructField
	t := reflect.TypeOf(ApkMetadata{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == "Dependencies" {
			f.Type = reflect.TypeOf(spaceDelimitedStringSlice{})
		}
		fields = append(fields, f)
	}
	apkMetadata := reflect.StructOf(fields)
	inst := reflect.New(apkMetadata)
	if err := json.Unmarshal(data, inst.Interface()); err != nil {
		return err
	}

	return mapstructure.Decode(inst.Elem().Interface(), m)
}

func (a *spaceDelimitedStringSlice) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}

	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return err
	}

	switch obj := jsonObj.(type) {
	case string:
		if obj == "" {
			*a = nil
		} else {
			*a = strings.Split(obj, " ")
		}
		return nil
	case []interface{}:
		s := make([]string, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return fmt.Errorf("invalid type for string array element: %T", v)
			}
			s = append(s, value)
		}
		*a = s
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("invalid type for string array: %T", obj)
	}
}

// ApkFileRecord represents a single file listing and metadata from a APK DB entry (which may have many of these file records).
type ApkFileRecord struct {
	Path        string       `json:"path"`
	OwnerUID    string       `json:"ownerUid,omitempty"`
	OwnerGID    string       `json:"ownerGid,omitempty"`
	Permissions string       `json:"permissions,omitempty"`
	Digest      *file.Digest `json:"digest,omitempty"`
}

func (m ApkMetadata) OwnedFiles() (result []string) {
	s := strset.New()
	for _, f := range m.Files {
		if f.Path != "" {
			s.Add(f.Path)
		}
	}
	result = s.List()
	sort.Strings(result)
	return result
}
