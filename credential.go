package main

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Credential struct {
  Orgs []string
  Token string
}

const CREDENTIAL_FILENAME = ".qp.credential.yaml"

func ReadCredential() (Credential, error) {
  var cred Credential
  var userHome, err = os.UserHomeDir()
  if err != nil {
    return Credential{}, err
  }

  fileContent, err := ioutil.ReadFile(path.Join(userHome, CREDENTIAL_FILENAME))
  if err != nil {
    return Credential{}, err
  }
  err = yaml.Unmarshal(fileContent, &cred)
  if err != nil {
    return Credential{}, err
  }
  return cred, nil
}

