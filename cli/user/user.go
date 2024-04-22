package user

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"webstack/metier/users"
)

type UserCfg struct {
	Email string `json:"email"`
	Mdp   string `json:"mdp"`
}

const configFilePath = "../.cfg/config.json"
const ERR_CFG = "erreur de chargement de la config:"
const ERR_SAVE_CFG = "erreur lors de l'enregistrement des informations:"
const ERR_WRITE = "error writing updated config:"
const ERR_READ = "error readfile config:"
const ERR_UNMARSH = "error unmarshal data:"
const ERR_MARSH = "error marshal updated data:"
const ERR_SIGNIN = "erreur pendant l'enregistrement d'un nouvel utilisateur :"

func Auth(f func(u users.User)) func(u users.User) {
	configData, err := LoadConfig()
	if err != nil {
		fmt.Println(ERR_CFG, err)
		return nil
	}
	if configData.Email != "" && configData.Mdp != "" {
		fmt.Println("Utilisateur connecté :", configData.Email)
		u, err := users.NewUser(configData.Email, configData.Mdp)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		f(u)
	} else {
		fmt.Print("Aucun utilisateur connecté: Voulez-vous vous connecter (c) ou vous inscrire (i)? ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			choice := strings.ToLower(scanner.Text())
			switch choice {
			case "c":
				u, _ := Login()
				f(u)
			case "i":
				u := Signin()
				f(u)
			default:
				fmt.Println("Choix invalide. Choisissez 'c' pour vous connecter ou 'i' pour vous inscrire.")
			}
		}
	}
	return f
}

func LoadConfig() (configUser UserCfg, err error) {
	if _, err = os.Stat(configFilePath); os.IsNotExist(err) {
		configUser = UserCfg{}
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

func SaveConfig(configUser UserCfg) error {
	fileContent, err := json.MarshalIndent(configUser, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, fileContent, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ClearUserConfig() error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		err = fmt.Errorf(" %v", err)
		return err
	}

	var u UserCfg
	err = json.Unmarshal(data, &u)
	if err != nil {
		err = fmt.Errorf("%v %v", ERR_UNMARSH, err)
		return err
	}

	u.Email = ""
	u.Mdp = ""

	updatedData, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		err = fmt.Errorf("%v %v", ERR_MARSH, err)
		return err
	}
	err = os.WriteFile(configFilePath, updatedData, 0644)
	if err != nil {
		err = fmt.Errorf("%v %v", ERR_WRITE, err)
		return err
	}

	fmt.Println("Utilisateur bien déconnecté !")
	return nil
}

func NewUserCfg(email string, mdp string) (u UserCfg) {
	u.Email = email
	u.Mdp = mdp
	return u
}

func Login() (u users.User, err error) {
	configData, err := LoadConfig()
	if err != nil {
		fmt.Println(ERR_CFG, err)
		return
	}
	fmt.Print("Entrez votre email: ")
	fmt.Scan(&configData.Email)
	fmt.Print("Entrez votre mot de passe: ")
	fmt.Scan(&configData.Mdp)

	u, err = users.Login(configData.Email, configData.Mdp)
	if err != nil {
		fmt.Println(err, "\nSi ce n'est pas déjà fait pensez à vous inscrire avec la commande signin !")
		ClearUserConfig()
		return u, err
	}
	err = SaveConfig(configData)
	if err != nil {
		fmt.Println(ERR_SAVE_CFG, err)
		return u, err
	}
	fmt.Println("Informations sauvegardées.")
	return u, nil
}

func Signin() (u users.User) {
	configData, err := LoadConfig()
	if err != nil {
		fmt.Println(ERR_CFG, err)
	}
	var confirmmdp string
	fmt.Print("Entrez votre email: ")
	fmt.Scan(&configData.Email)
	fmt.Print("Choisissez un mot de passe: ")
	fmt.Scan(&configData.Mdp)
	fmt.Print("Confirmez votre mot de passe: ")
	fmt.Scan(&confirmmdp)
	u, err = users.Signin(configData.Email, configData.Mdp, confirmmdp)
	if err != nil {
		fmt.Println(ERR_SIGNIN, err)
		return
	}
	err = SaveConfig(configData)
	if err != nil {
		fmt.Println(ERR_SAVE_CFG, err)
		return
	}
	fmt.Println("Informations sauvegardées.")
	return u
}
