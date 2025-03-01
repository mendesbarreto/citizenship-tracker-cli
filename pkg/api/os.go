package api

import "os"

func GetUserUci() string {
	uci := os.Getenv("TRACKER_UCI")

	if uci == "" {
		credentialsKeychain, err := GetFromKeychain()
		if err != nil {
			return ""
		}

		uci = credentialsKeychain.Account
	}

	return uci
}

func GetUserPassword() string {
	password := os.Getenv("TRACKER_PASSWORD")

	if password == "" {
		credentialsKeychain, err := GetFromKeychain()
		if err != nil {
			return ""
		}

		password = credentialsKeychain.Password
	}

	return password
}

func GetUserApplicationNumber() string {
	applicationNumber := os.Getenv("TRACKER_APPLICATION_NUMBER")

	if applicationNumber == "" {
		credentialsKeychain, err := GetFromKeychain()
		if err != nil {
			return ""
		}

		applicationNumber = credentialsKeychain.ApplicationNumber
	}

	return applicationNumber
}
