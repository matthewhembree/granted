package registry

import (
	"bytes"
	"regexp"
	"text/template"

	grantedConfig "github.com/common-fate/granted/pkg/config"
)

const AUTO_GENERATED_MSG string = `# Granted-Registry Autogenerated Section. DO NOT EDIT.
# This section is automatically generated by Granted (https://granted.dev). Manual edits to this section will be overwritten.
# To edit, clone your profile registry repo, edit granted.yml, and push your changes. You may need to make a pull request depending on the repository settings.
# To stop syncing and remove this section, run 'granted registry remove'.`

func getAutogeneratedTemplate() string {
	return AUTO_GENERATED_MSG
}

func interpolateVariables(r *Registry, value string, profileName string) (string, error) {
	var variables = r.Variables

	for k, v := range r.Variables {
		variables[k] = v
	}

	gConf, err := grantedConfig.Load()
	if err != nil {
		return "", err
	}

	for k, v := range gConf.ProfileRegistry.RequiredKeys {
		variables[k] = v
	}

	// profile is automatically figured out by granted.
	variables["profile"] = profileName

	tmpl, err := template.New("registry-variable").Parse(value)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, variables)
	if err != nil {
		panic(err)
	}
	return buf.String(), nil
}

func containsTemplate(text string) bool {
	re := regexp.MustCompile(`{{\s+(\w+)?.\w+\s+}}`)

	return re.MatchString(text)
}
