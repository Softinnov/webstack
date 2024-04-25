package user

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"webstack/metier/users"
)

type UCfg struct {
	Email string `json:"email"`
	Mdp   string `json:"mdp"`
}

const CfgFilePath = "./.cfg/config.json"
const ErrCfg = "erreur de chargement de la config"
const ErrSaveCfg = "erreur lors de l'enregistrement des informations"
const ErrWrite = "erreur writing updated config"
const ErrRead = "erreur readfile config"
const ErrUnmarsh = "erreur unmarshal data"
const ErrMarsh = "erreur marshal updated data"
const ErrSignin = "erreur pendant l'enregistrement d'un nouvel utilisateur"
const SavedInfo = "Informations sauvegardées."
const NotSignedin = "\nSi ce n'est pas déjà fait pensez à vous inscrire avec la commande signin !"
const NouserCfg = "Aucun utilisateur connecté: Voulez-vous vous connecter (c) ou vous inscrire (i)? "
const InvChoice = "Choix invalide. Choisissez 'c' pour vous connecter ou 'i' pour vous inscrire."
const ErrScan = "erreur de scan"
const FormativeStr = "%v : %v"

func Auth(f func(u users.User), configFilePath string) func(u users.User) {
	configData, err := LoadConfig(configFilePath)
	if err != nil {
		fmt.Println(ErrCfg, err)
		return nil
	}

	if configData.Email != "" && configData.Mdp != "" {
		u, err := users.NewUser(configData.Email, configData.Mdp)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		f(u)
	} else {
		fmt.Print(NouserCfg)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			choice := strings.ToLower(scanner.Text())
			switch choice {
			case "c":
				u, err := Login(configFilePath)
				if err != nil {
					fmt.Println(err)
				}
				f(u)
			case "i":
				u, err := Signin(configFilePath)
				if err != nil {
					fmt.Println(err)
				}
				f(u)
			default:
				fmt.Println(InvChoice)
				return nil
			}
		}
	}

	return f
}

func LoadConfig(configFilePath string) (configUser UCfg, err error) {
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		configUser = UCfg{}
	} else {
		fileContent, err := os.ReadFile(configFilePath)
		if err != nil {
			return configUser, err
		}
		err = json.Unmarshal(fileContent, &configUser)
		if err != nil {
			return configUser, err
		}
	}

	return configUser, nil
}

func SaveConfig(configUser UCfg, configFilePath string) error {
	fileContent, err := json.MarshalIndent(configUser, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, fileContent, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func ClearUserConfig(configFilePath string) error {
	var u UCfg

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return fmt.Errorf(FormativeStr, ErrUnmarsh, err)
	}

	u.Email = ""
	u.Mdp = ""

	updatedData, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf(FormativeStr, ErrMarsh, err)
	}

	err = os.WriteFile(configFilePath, updatedData, 0o644)
	if err != nil {
		return fmt.Errorf(FormativeStr, ErrWrite, err)
	}

	return nil
}

func NewUserCfg(email, mdp string) (u UCfg) {
	u.Email = email
	u.Mdp = mdp

	return u
}

func Login(configFilePath string) (u users.User, err error) {
	configData, err := LoadConfig(configFilePath)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrCfg, err)
	}

	fmt.Print("Entrez votre email: ")

	_, err = fmt.Scan(&configData.Email)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrScan, err)
	}

	fmt.Print("Entrez votre mot de passe: ")

	_, err = fmt.Scan(&configData.Mdp)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrScan, err)
	}

	u, err = users.Login(configData.Email, configData.Mdp)
	if err != nil {
		fmt.Println(NotSignedin)

		err2 := ClearUserConfig(configFilePath)
		if err2 != nil {
			return u, err2
		}

		return u, err
	}

	err = SaveConfig(configData, configFilePath)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrSaveCfg, err)
	}

	fmt.Println(SavedInfo)

	return u, nil
}

func Signin(configFilePath string) (u users.User, err error) {
	var confirmmdp string

	configData, err := LoadConfig(configFilePath)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrCfg, err)
	}

	fmt.Print("Entrez votre email: ")

	_, err = fmt.Scan(&configData.Email)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrScan, err)
	}

	fmt.Print("Choisissez un mot de passe: ")

	_, err = fmt.Scan(&configData.Mdp)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrScan, err)
	}

	fmt.Print("Confirmez votre mot de passe: ")

	_, err = fmt.Scan(&confirmmdp)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrScan, err)
	}

	u, err = users.Signin(configData.Email, configData.Mdp, confirmmdp)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrSignin, err)
	}

	err = SaveConfig(configData, configFilePath)
	if err != nil {
		return u, fmt.Errorf(FormativeStr, ErrSaveCfg, err)
	}

	fmt.Println(SavedInfo)

	return u, nil
}
