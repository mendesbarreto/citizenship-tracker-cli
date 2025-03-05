package api

import (
	"fmt"

	"github.com/keybase/go-keychain"
)

const (
	ServiceKey string = "citizenship-tracker-cli"
	BundleID   string = "com.douglas.mendes.citizenship-tracker-cli"
)

type KeychainItem struct {
	Account           string
	Service           string
	ApplicationNumber string
	Password          string
}

func AddToKeychain(account string, password string, label string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(ServiceKey)
	item.SetAccount(account)
	item.SetLabel(label)
	item.SetAccessGroup(BundleID)
	item.SetData([]byte(password))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	return keychain.AddItem(item)
}

func GetFromKeychain() (KeychainItem, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(ServiceKey)
	query.SetAccessGroup(BundleID)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	query.SetReturnAttributes(true)
	results, err := keychain.QueryItem(query)

	if err != nil {
		return KeychainItem{}, err
	} else if len(results) != 1 {
		fmt.Println("[Keychain] No results found")
		return KeychainItem{}, fmt.Errorf("no results found for this keychain service: %s", ServiceKey)
	}

	firstResult := results[0]

	// fmt.Printf("[Keychain] Account: %s\n", firstResult.Account)
	// fmt.Printf("[Keychain] Service: %s\n", firstResult.Service)
	// fmt.Printf("[Keychain] Application Number: %s\n", firstResult.Label)
	// fmt.Printf("[Keychain] Password: %s\n", string(firstResult.Data))

	return KeychainItem{
		Account:           firstResult.Account,
		Service:           firstResult.Service,
		ApplicationNumber: firstResult.Label,
		Password:          string(firstResult.Data),
	}, nil
}

func ExistsOnKeychain() bool {
	item, err := GetFromKeychain()
	return err == nil && item.Service == ServiceKey
}
