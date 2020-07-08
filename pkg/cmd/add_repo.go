package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/prompt"
	"github.com/ZupIT/ritchie-cli/pkg/stdin"
)

type AddRepoCmd struct {
	client *http.Client
	repo   formula.RepositoryAdder
	prompt.InputText
	prompt.InputPassword
	prompt.InputURL
	prompt.InputList
	prompt.InputBool
	prompt.InputInt
}

func NewAddRepoCmd(
	client *http.Client,
	repo formula.RepositoryAdder,
	inText prompt.InputText,
	inPass prompt.InputPassword,
	inUrl prompt.InputURL,
	inList prompt.InputList,
	inBool prompt.InputBool,
	inInt prompt.InputInt,
) *cobra.Command {
	addRepo := AddRepoCmd{
		client:        client,
		repo:          repo,
		InputText:     inText,
		InputURL:      inUrl,
		InputList:     inList,
		InputBool:     inBool,
		InputInt:      inInt,
		InputPassword: inPass,
	}
	cmd := &cobra.Command{
		Use:     "repo",
		Short:   "Add a repository.",
		Example: "rit add repo ",
		RunE:    RunFuncE(addRepo.runStdin(), addRepo.runPrompt()),
	}
	cmd.LocalFlags()

	return cmd
}

func (a AddRepoCmd) runPrompt() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {
		name, err := a.Text("Repository name: ", true)
		if err != nil {
			return err
		}

		isPrivate, err := a.Bool("Is a private repository? ", []string{"no", "yes"})
		if err != nil {
			return err
		}

		var token string
		if isPrivate {
			token, err = a.Password("Personal access tokens: ")
			if err != nil {
				return err
			}
		}

		url, err := a.URL("Repository URL: ", "https://github.com/kaduartur/ritchie-formulas")
		if err != nil {
			return err
		}

		tags, err := a.tags(url, token)
		if err != nil {
			return err
		}

		var tagNames []string
		for k := range tags {
			tagNames = append(tagNames, k)
		}

		version, err := a.List("Select a tag version: ", tagNames)
		if err != nil {
			return err
		}

		zipUrl := tags[version]

		priority, err := a.Int("Set the priority [ps.: 0 is higher priority, the lower higher the priority] :")
		if err != nil {
			return err
		}

		repository := formula.Repo{
			Name:     name,
			Token:    token,
			ZipUrl:   zipUrl,
			Version:  version,
			Priority: int(priority),
		}

		if err := a.repo.Add(repository); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was added with success, now you can use your formulas with the Ritchie!", repository.Name)
		prompt.Success(successMsg)
		return nil
	}
}

func (a AddRepoCmd) runStdin() CommandRunnerFunc {
	return func(cmd *cobra.Command, args []string) error {

		r := formula.Repo{}

		err := stdin.ReadJson(os.Stdin, &r)
		if err != nil {
			prompt.Error(stdin.MsgInvalidInput)
			return err
		}

		if err := a.repo.Add(r); err != nil {
			return err
		}

		successMsg := fmt.Sprintf("The %q repository was added with success, now you can use your formulas with the Ritchie!", r.Name)
		prompt.Success(successMsg)
		return nil
	}
}

func (a AddRepoCmd) tags(url string, token string) (formula.Tags, error) {
	apiUrl, err := tagsUrl(url)
	if err != nil {
		return formula.Tags{}, err
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return formula.Tags{}, err
	}

	if token != "" {
		authToken := fmt.Sprintf("token %s", token)
		req.Header.Add(headers.Authorization, authToken)
	}

	req.Header.Add(headers.Accept, "application/vnd.github.v3+json")
	resp, err := a.client.Do(req)
	if err != nil {
		return formula.Tags{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return formula.Tags{}, err
		}
		return formula.Tags{}, errors.New(string(b))
	}

	var tags []formula.Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return formula.Tags{}, err
	}

	tagsUrl := make(formula.Tags)
	for _, tag := range tags {
		tagsUrl[tag.Name] = tag.ZipUrl
	}

	return tagsUrl, nil
}

func tagsUrl(url string) (string, error) {
	split := strings.Split(url, "/")
	repo := split[len(split)-1]
	owner := split[len(split)-2]

	return fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo), nil
}
