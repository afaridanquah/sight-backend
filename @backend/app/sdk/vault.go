package sdk

import (
	"os"

	"bitbucket.org/msafaridanquah/verifylab-service/foundation/envvar/vault"
	"bitbucket.org/msafaridanquah/verifylab-service/foundation/ierr"
)

// NewVaultProvider instantiates the Vault client using configuration defined in environment variables.
func NewVaultProvider() (*vault.Provider, error) {
	// XXX: We will revisit this code in future episodes replacing it with another solution
	vaultPath := os.Getenv("VAULT_PATH")
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	// XXX: -

	provider, err := vault.New(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		return nil, ierr.WrapErrorf(err, ierr.ErrorCodeUnknown, "vault.New")
	}

	return provider, nil
}
