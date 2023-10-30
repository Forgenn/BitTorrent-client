package main

import (
	"os"
	"strconv"
	"strings"
)

func decodeTorrentFile(filePath string) interface{} {
	var fileContents, err = os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return decodeBencode(string(fileContents))
}

func decodeBencode(bencode string) interface{} {
	var result, _ = decodeBencodeUtil(bencode, 0)
	return result
}

func decodeBencodeUtil(bencode string, length int) (interface{}, int) {
	var i int = 0

	switch firstChar := bencode[i]; firstChar {
	// Decode an integer
	case 'i':
		var delimiterIndex = strings.Index(bencode[i:], "e")
		var resultTmpInt, _ = strconv.Atoi(bencode[i+1 : delimiterIndex])
		return resultTmpInt, i + delimiterIndex + 1
	// Decode a list
	case 'l':
		i++
		var result []interface{}
		var resultTmp interface{}
		for bencode[i] != 'e' {
			resultTmp, length = decodeBencodeUtil(bencode[i:], length)
			result = append(result, resultTmp)
			i += length
		}
		return result, i + 1
	// Decode dictionary
	case 'd':
		i++
		var value interface{}
		dict := make(map[string]interface{})
		for bencode[i] != 'e' {
			var key, lengthKey = decodeBencodeUtil(bencode[i:], length)
			i += lengthKey
			value, length = decodeBencodeUtil(bencode[i:], length+lengthKey)
			i += length
			dict[key.(string)] = value
		}
		return dict, i + 1
	case 'e':
		break

	// Decode Byte Strings
	default:
		var colonIndex = strings.Index(bencode[i:], ":")
		var lengthInt, _ = strconv.Atoi(bencode[i : i+colonIndex])
		return bencode[i+colonIndex+1 : i+colonIndex+lengthInt+1], i + colonIndex + lengthInt + 1
	}
	return "error", i
}
