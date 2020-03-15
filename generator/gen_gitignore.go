package generator

import "io/ioutil"

func (gen *caGen) GenGitIgnore(dirName string) error {
	gitignore := []byte(`vendor
.env
.DS_Store`)

	err := ioutil.WriteFile("./"+dirName+"/.gitignore", gitignore, 0644)
	if err != nil {
		return err
	}

	return nil
}
