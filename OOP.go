package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type DataSet struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

type Command interface {
	Name() string
	ParseFlags(args []string) error
	Execute() error
}

type CreateCommand struct {
	name    string
	user    string
	project string
	fs      *flag.FlagSet
}
type GetCommand struct {
	name string
	fs   *flag.FlagSet
}
type ListCommand struct {
	fs   *flag.FlagSet
}
type DeleteCommand struct {
	name string
	fs   *flag.FlagSet
}
// type HelpCommand struct{
// 	fs *flag.FlagSet
// }

func printWelcome() {
	fmt.Println("Добро поджаловать в gogomts, CLI инструмент для управления профилями")
	fmt.Println("Для начала пропишите'./gogomts --help' для большей информации о возможностях")
}

func printHelp() {
	fmt.Println(`Доступные комманды:
  profile create --name=<name> --user=<user> --project=<project> - Создание нового профиля
  profile get --name=<name>                                      - Вывидение информации о профиле
  profile list                                                   - Список профилей
  profile delete --name=<name>                                   - Удалить профиль

	Пример использования:
  ./gogomts profile create --name=test --user=example --project=new-project`)
}

// func (c *HelpCommand) Name() string {
// 	return "help"
// }

// func (c *HelpCommand) ParseFlags(args []string) error {
// 	c.fs = flag.NewFlagSet("help", flag.ExitOnError)
// 	return c.fs.Parse(args)
// }

// func (c *HelpCommand) Execute() error{
// 	fmt.Println(`Доступные комманды:
//   profile create --name=<name> --user=<user> --project=<project> - Создание нового профиля
//   profile get --name=<name>                                      - Вывидение информации о профиле
//   profile list                                                   - Список профилей
//   profile delete --name=<name>                                   - Удалить профиль

// 	Пример использования:
//   ./gogomts profile create --name=test --user=example --project=new-project`)
// fmt.Println("Доступные профили:")
//   return nil
// }


func (c *CreateCommand) Name() string {
	return "create"
}

func (c *CreateCommand) ParseFlags(args []string) error {
	c.fs = flag.NewFlagSet("create", flag.ExitOnError)
	c.fs.StringVar(&c.name, "name", "", "Profile name")
	c.fs.StringVar(&c.user, "user", "", "User name")
	c.fs.StringVar(&c.project, "project", "", "Project name")
	return c.fs.Parse(args)
}

func (c *CreateCommand) Execute() error {
	userInfo := DataSet{
		User:    c.user,
		Project: c.project,
	}
	data, err := yaml.Marshal(&userInfo)
	if err != nil {
		return fmt.Errorf("ошибка при создании YAML: %v", err)
	}
	err = os.WriteFile(c.name+".yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}
	return nil
}

func (c *GetCommand) Name() string {
	return "get"
}

func (c *GetCommand) ParseFlags(args []string) error {
	c.fs = flag.NewFlagSet("get", flag.ExitOnError)
	c.fs.StringVar(&c.name, "name", "", "Profile name")
	return c.fs.Parse(args)
}

func (c *GetCommand) Execute() error {
	if c.name == "" {
		return fmt.Errorf("необходим флаг --name")
	}
	data, err := os.ReadFile(c.name + ".yaml")
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}
	var profile DataSet
	err = yaml.Unmarshal(data, &profile)
	if err != nil {
		return fmt.Errorf("ошибка загрузки YAML: %v", err)
	}
	fmt.Printf("Profile: %s\nUser: %s\nProject: %s\n", c.name, profile.User, profile.Project)
	return nil
}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) ParseFlags(args []string) error {
	c.fs = flag.NewFlagSet("list", flag.ExitOnError)
	return c.fs.Parse(args)
}

func (c *ListCommand) Execute() error {
	files, err := filepath.Glob("*.yaml")
	if err != nil {
		return fmt.Errorf("ошибка создания списка профилей: %v", err)
	}
	if len(files) == 0 {
		fmt.Println("Нет созданных профилей")
		return nil
	}
	fmt.Println("Доступные профили:")
	for _, file := range files {
		profileName := file[:len(file)-5]
		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("ошибка чтения профиля %s: %v", profileName, err)
			continue
		}
		var profile DataSet
		err = yaml.Unmarshal(data, &profile)
		if err != nil {
			log.Printf("ошибка загрузки профилей %s: %v", profileName, err)
			continue
		}
		fmt.Printf("- %s (User: %s, Project: %s)\n", profileName, profile.User, profile.Project)
	}
	return nil
}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func (c *DeleteCommand) ParseFlags(args []string) error {
	c.fs = flag.NewFlagSet("delete", flag.ExitOnError)
	c.fs.StringVar(&c.name, "name", "", "Profile name")
	return c.fs.Parse(args)
}

func (c *DeleteCommand) Execute() error {
	if c.name == "" {
		return fmt.Errorf("необходим флаг --name")
	}
	err := os.Remove(c.name + ".yaml")
	if err != nil {
		return fmt.Errorf("ошибка удаления профиля: %v", err)
	}
	return nil
}

func main() {
	rootFS := flag.NewFlagSet("gogomts", flag.ExitOnError)
	rootHelp := rootFS.Bool("help", false, "Show welcome message")

	if *rootHelp {
		printHelp()
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		printWelcome()
		os.Exit(1)
	}
	commands := map[string]Command{
		"create": &CreateCommand{},
		"get":    &GetCommand{},
		"list":   &ListCommand{},
		"delete": &DeleteCommand{},
		// "help":   &HelpCommand{},
	}

	if os.Args[1] == "help" {
		printHelp()
		os.Exit(1)
	}

	if os.Args[1] != "profile"{
	 		fmt.Printf("неизвестная командп: %s\n", os.Args[1])
	 		printHelp()
	 		os.Exit(1)
	 	}

	if len(os.Args) < 3 {
		fmt.Println("Пожалуйста укажите действие с профилем")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]
	cmd, exists := commands[subcommand]
	if !exists {
		fmt.Printf("Незвестное действие: %s\n", subcommand)
		printHelp()
		os.Exit(1)
	}
	if err := cmd.ParseFlags(os.Args[3:]); err != nil {
		log.Fatalf("Ошибка парсинга флагов: %v", err)
	}
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Ошибка исполнения команды: %v", err)
	}
}