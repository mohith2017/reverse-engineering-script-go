package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"github.com/vmihailenco/msgpack/v5"
)

func removeNonHex(input string) string {
	nonHexPattern := regexp.MustCompile("[^0-9a-fA-F]")
	cleanedString := nonHexPattern.ReplaceAllString(input, "")

	return cleanedString
}

func decryptAscii(input string, change int) string {
	var decrypted string

	for _, char := range input {
		decryptedChar := int(char) - change
		decrypted += string(decryptedChar)
	}

	return decrypted
}

func extractNumber(input string) string {
	var changeString string

	for i := 1; i < len(input); i++ {
		if input[i] == ' ' {
			for j := i+1; j < len(input) && input[j] != ' '; j++ {
				changeString += string(input[j])
			}
			break
		}
	}

	return changeString
}

func customHexDecode(encodedString string, customCharset string) (string, error) {

	standardHexSet := "0123456789abcdef"
	decodeMap := make(map[rune]rune)

	encodedString = strings.ReplaceAll(encodedString, " ", "")
	customCharset = strings.ReplaceAll(customCharset, " ", "")

	for i, char := range []rune(customCharset) {
		decodeMap[char] = []rune(standardHexSet)[i]
	}

	
	// Decode each character in the encoded string
	var decodedRunes []rune
	for _, char := range []rune(encodedString) {
		// fmt.Println("Character: ", char)
		decodedChar, ok := decodeMap[char]
		if !ok {
			return "", fmt.Errorf("invalid character in encoded string: %c", char)
		}
		decodedRunes = append(decodedRunes, decodedChar)
	}

	return string(decodedRunes), nil
}


func unscramble(encryptedString string, encodedMessagePack string) (string, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(encodedMessagePack)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return "", err
	}

	var positions []int
	err = msgpack.Unmarshal(decodedBytes, &positions)
	if err != nil {
		fmt.Println("Error decoding MessagePack:", err)
		return "", err
	}

	// Create a new byte slice
	unscrambled := make([]byte, len(encryptedString))

	for i, pos := range positions {
		unscrambled[pos] = encryptedString[i]
	}

	decryptedString := string(unscrambled)
	fmt.Println("Decrypted string:", decryptedString)

	return decryptedString, nil
}

func getRequest(challengeURL string) (map[string]interface{}, error) {
	challenge, err := http.Get(challengeURL)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return nil, err
	}
	defer challenge.Body.Close()

	
	body, err := ioutil.ReadAll(challenge.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return nil, err
		}

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	fmt.Println("Response:", string(body))
	return response, nil
}

func main() {
	var nextChallengeURL string
	var challengeURL string
	var response map[string]interface{}

	//Base Level
	applicationEmail := "mohith.ny2024@gmail.com"
	baseURL := "https://ciphersprint.pulley.com/"
	challengeURL = baseURL + applicationEmail

	response, err := getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level  0
	nextChallengeURL = response["encrypted_path"].(string)
	challengeURL = baseURL + nextChallengeURL
	fmt.Println("Next challenge URL:", challengeURL)

	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level  1
	nextChallengeURL = response["encrypted_path"].(string)
	decodedPath, err := base64.StdEncoding.DecodeString(nextChallengeURL[5:])
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}
	decodedString := string(decodedPath)
	challengeURL = baseURL + "task_" + decodedString
	fmt.Println("Next challenge URL:", challengeURL)
	
	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level  2
	nextChallengeURL = response["encrypted_path"].(string)
	fmt.Println("Next challenge URL:", nextChallengeURL[5:])
	cleanedString := removeNonHex(nextChallengeURL[5:])
	challengeURL = baseURL + "task_" + cleanedString
	fmt.Println("Next challenge URL:", challengeURL)
	
	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level 3
	var changeString string
	nextChallengeURL = response["encrypted_path"].(string)

	encryptionMethod := response["encryption_method"].(string)
	changeString = extractNumber(encryptionMethod)

	change, err := strconv.Atoi(changeString)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Change in encryption: ", change)
	decryptedString := decryptAscii(nextChallengeURL[5:], change)
	fmt.Println("Decrypted string: ", decryptedString)
	challengeURL = baseURL + "task_" + decryptedString
	fmt.Println("Next challenge URL:", challengeURL)
	
	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level 4
	nextChallengeURL = response["encrypted_path"].(string)
	customCharset := (response["encryption_method"].(string))[41:]
	fmt.Println("Custom Hex Char set: ", customCharset)
	customDecodedString, err := customHexDecode(nextChallengeURL[5:], customCharset)
	if err != nil {
		fmt.Println("Error in decoding the Hex string: ", err)
		return
	}

	fmt.Println("Decoded Hex string: ", customDecodedString)
	challengeURL = baseURL + "task_" + customDecodedString
	fmt.Println("Next challenge URL:", challengeURL)
	
	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//Level 5
	nextChallengeURL = response["encryption_method"].(string)
	fmt.Println("Postion: ", nextChallengeURL[61:])
	unscrambledString, err := unscramble((response["encrypted_path"].(string))[5:], nextChallengeURL[61:])
	if err != nil {
		fmt.Println("Error in unscrambling data: ", err)
		return
	}

	challengeURL = baseURL + "task_" + unscrambledString
	fmt.Println("Next challenge URL:", challengeURL)
	
	response, err = getRequest(challengeURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
