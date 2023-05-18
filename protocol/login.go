package protocol

import (
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PondWader/GoPractice/database"
	"github.com/PondWader/GoPractice/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginStartPacket struct {
	Name string `type:"String"`
}

type EncryptionRequestPacket struct {
	ServerID          string `type:"String"`
	PublicKeyLength   int    `type:"VarInt"`
	PublicKey         []byte `type:"ByteArray"`
	VerifyTokenLength int    `type:"VarInt"`
	VerifyToken       []byte `type:"ByteArray"`
}

type EncryptionResponsePacket struct {
	SharedSecretLength int    `type:"VarInt"`
	SharedSecret       []byte `type:"ByteArray" length:"SharedSecretLength"`
	VerifyTokenLength  int    `type:"VarInt"`
	VerifyToken        []byte `type:"ByteArray" length:"VerifyTokenLength"`
}

type LoginSuccessPacket struct {
	UUID     string `type:"String"`
	Username string `type:"String"`
}

type SessionServerResp struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature"`
	} `json:"properties"`
}

type SetCompressionPacket struct {
	Threshold int `type:"VarInt"`
}

func (client *ProtocolClient) login() {
	client.state = "login"

	packetId, data, err := client.readPacket()
	if err != nil {
		return
	}
	if packetId != 0 {
		client.Disconnect("Bad packet ID")
		return
	}

	packet := &LoginStartPacket{}
	err = client.deserialize(data, packet)
	if err != nil {
		client.Disconnect(err.Error())
		return
	}
	client.Username = packet.Name

	pubKey, _ := x509.MarshalPKIXPublicKey(&client.server.GetKeyPair().PublicKey)

	verifyToken := make([]byte, 4)
	rand.Reader.Read(verifyToken)

	if err := client.WritePacket(0x01, Serialize(&EncryptionRequestPacket{
		PublicKeyLength:   len(pubKey),
		PublicKey:         pubKey,
		VerifyTokenLength: 4,
		VerifyToken:       verifyToken,
	})); err != nil {
		return
	}

	packetId, data, err = client.readPacket()
	if err != nil {
		return
	}
	if packetId != 0x01 {
		client.Disconnect("Bad packet ID")
		return
	}

	encryptionResponse := &EncryptionResponsePacket{}
	if err := client.deserialize(data, encryptionResponse); err != nil {
		return
	}

	if len(encryptionResponse.VerifyToken) != 128 {
		client.Disconnect("Invalid verify token length!")
		return
	}
	decryptedVerifyToken, err := client.server.GetKeyPair().Decrypt(rand.Reader, encryptionResponse.VerifyToken, nil)
	if err != nil {
		client.Disconnect(err.Error())
		return
	}
	if !compareVerifyTokens(decryptedVerifyToken, verifyToken) {
		client.Disconnect("Verify tokens do not match!")
		return
	}

	if len(encryptionResponse.SharedSecret) != 128 {
		client.Disconnect("Invalid shared secret length!")
		return
	}

	sharedSecret, err := client.server.GetKeyPair().Decrypt(rand.Reader, encryptionResponse.SharedSecret, nil)
	if err != nil {
		client.Disconnect(err.Error())
		return
	}

	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		client.Disconnect(err.Error())
		return
	}
	client.decrypter = utils.NewCFB8Decrypter(block, sharedSecret)
	client.encrypter = utils.NewCFB8Encrypter(block, sharedSecret)
	client.encryption = true

	if client.server.GetConfig().OnlineMode == true {
		db := client.server.GetDatabase()

		hash := sha1.New()
		hash.Write(append(sharedSecret, pubKey...))
		hashHex := hex.EncodeToString(hash.Sum(nil))

		resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + client.Username + "&serverId=" + hashHex)
		if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 204) {
			client.Disconnect("The auth servers appear to be down.")
			utils.Info(resp.StatusCode)
			return
		} else if resp.StatusCode == 204 {
			// If status 204 is returned, the player info should be cached
			var result database.SessionCache
			tx := db.First(&result, "username = ?", client.Username)
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				client.Disconnect("Authentication temporarily unavailable")
				return
			} else if tx.Error != nil {
				utils.Error("Error retrieving session cache record:", tx.Error)
				client.Disconnect("Error communicating with database")
				return
			}
			uid, _ := uuid.Parse(result.Uuid)
			client.Uuid = &uid
			client.Skin = result.Textures
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			userInfo := &SessionServerResp{}
			if err := json.Unmarshal(body, userInfo); err != nil {
				client.Disconnect("Received an invalid response from the session server.")
				return
			}

			db.Save(&database.SessionCache{
				Uuid:       userInfo.Id,
				Username:   userInfo.Name,
				Textures:   userInfo.Properties[0].Value,
				TimeCached: time.Now(),
			})

			uid, _ := uuid.Parse(userInfo.Id)
			client.Uuid = &uid
			client.Skin = userInfo.Properties[0].Value
		}
	} else {
		client.Uuid = createOfflineUUID("OfflinePlayer:" + client.Username)
	}

	utils.Info(client.Username, "has succesfully logged in.")

	if err := client.WritePacket(0x03, Serialize(&SetCompressionPacket{
		Threshold: client.server.GetConfig().CompressionThreshold,
	})); err != nil {
		return
	}

	client.compression = true

	if err := client.WritePacket(0x02, Serialize(&LoginSuccessPacket{
		UUID:     client.Uuid.String(),
		Username: client.Username,
	})); err != nil {
		return
	}

	client.play()
}

func compareVerifyTokens(decryptedVerifyToken []byte, verifyToken []byte) bool {
	if len(decryptedVerifyToken) != len(verifyToken) {
		return false
	}

	for i := 0; i < len(verifyToken); i++ {
		if verifyToken[i] != decryptedVerifyToken[i] {
			return false
		}
	}

	return true
}

// https://github.com/openjdk-mirror/jdk7u-jdk/blob/f4d80957e89a19a29bb9f9807d2a28351ed7f7df/src/share/classes/java/util/UUID.java#L163
func createOfflineUUID(name string) *uuid.UUID {
	md5Hash := md5.Sum([]byte(name))
	md5Bytes := md5Hash[:]
	md5Bytes[6] &= 0x0f // clear version
	md5Bytes[6] |= 0x30 // set to version 3
	md5Bytes[8] &= 0x3f // clear variant
	md5Bytes[8] |= 0x80 // set to IETF variant
	uuidBytes, err := hex.DecodeString(fmt.Sprintf("%x", md5Bytes))
	if err != nil {
		panic(err)
	}
	u, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		panic(err)
	}
	return &u
}
