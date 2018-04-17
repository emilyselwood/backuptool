package conf

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"regexp"
)

/*
Config contains all the information about a backup task
*/
type Config struct {
	Dirs    []Dir  `json:"dirs,omitempty"`
	Version string `json:"version,omitempty"`
}

/*
Dir is a directory that should be backed up
*/
type Dir struct {
	Path     string `json:"path"`
	Include  string `json:"include"`
	Exclude  string `json:"exclude"`
	incRegex *regexp.Regexp
	excRegex *regexp.Regexp
}

/*
String returns a string representation for printing of the dir entry
*/
func (d *Dir) String() string {
	return fmt.Sprintf("%s\t%s\t%s", d.Path, d.Include, d.Exclude)
}

/*
ShouldInclude returns true if the given path should be included according to this rules
*/
func (d *Dir) ShouldInclude(path string) bool {

	d.compile()
	if d.incRegex != nil || d.excRegex != nil {

		if d.incRegex.MatchString(path) {
			return true
		}
		if d.excRegex.MatchString(path) {
			return false
		}

	}
	return true

}

func (d *Dir) compile() {
	if d.incRegex == nil && d.excRegex == nil {
		if d.Include != "" || d.Exclude != "" {
			var err error
			d.incRegex, err = regexp.Compile(d.Include)
			if err != nil {
				panic("could not compile include regex for dir " + d.Path + " " + err.Error())
			}
			d.excRegex, err = regexp.Compile(d.Exclude)
			if err != nil {
				panic("could not compile exclude regex for dir " + d.Path + " " + err.Error())
			}
		}
	}
}

/*
Default contains the very basic configuration
*/
var Default = Config{Dirs: []Dir{}, Version: "0.0.1"}

/*
ReadConfig loads the config file. There is no optional path here. It is always fixed to make
*/
func ReadConfig() (*Config, error) {

	p, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	var result Config
	dec := json.NewDecoder(f)

	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

/*
ReadConfigZip loads up a configuration from a zip file.
*/
func ReadConfigZip(f *zip.ReadCloser) (*Config, error) {

	cFile, err := findFileInZip(f, "backup.conf")
	if err != nil {
		return nil, err
	}

	c, err := cFile.Open()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	var result Config
	dec := json.NewDecoder(c)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func findFileInZip(root *zip.ReadCloser, name string) (*zip.File, error) {
	for _, e := range root.File {
		if e.Name == name {
			return e, nil
		}
	}
	return nil, fmt.Errorf("could not find file %s in zip", name)
}

/*
WriteConfig stores an updated configuration entry
*/
func WriteConfig(c *Config) error {
	p, err := getConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteConfigZip(c, f)
}

/*
WriteConfigZip writes the config to the given writer
*/
func WriteConfigZip(c *Config, w io.Writer) error {
	dec := json.NewEncoder(w)
	dec.SetIndent("", "  ")

	if err := dec.Encode(&c); err != nil {
		return err
	}

	return nil
}

func getConfigPath() (string, error) {
	home, err := getHomeDir()
	if err != nil {
		return "", err
	}

	p := path.Join(home, ".config/backup.conf")
	return p, nil
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
