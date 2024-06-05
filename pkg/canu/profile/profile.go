package profile

// reloging workflow from: https://github.com/aws/aws-sdk-go-v2/issues/1222

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/pkg/browser"
)

const (
	ssoCachePath = "sso/cache"
)

type profile struct {
	name       string
	configPath string
}

type tokenFields struct {
	StartURL              string `json:"startUrl,omitempty"`
	Region                string `json:"region,omitempty"`
	AccessToken           string `json:"accessToken,omitempty"`
	ExpiresAt             string `json:"expiresAt,omitempty"`
	ClientID              string `json:"clientId,omitempty"`
	ClientSecret          string `json:"clientSecret,omitempty"`
	RegistrationExpiresAt string `json:"registrationExpiresAt,omitempty"`
}

func New(name, configPath string) *profile {
	return &profile{
		name:       name,
		configPath: configPath,
	}
}

func (p *profile) CheckAWSProfileStatus() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(p.name))

	if err != nil {
		return fmt.Errorf("error while preparing aws config for profile %s: %v", p.name, err)
	}

	sharedConfig := getSharedConfig(&cfg)

	if err := validateSharedConfig(sharedConfig); err != nil {
		return err
	}

	// to check if we need to relog we just do `aws sts --get-caller-identity`
	stsClient := sts.NewFromConfig(cfg)
	output, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	isRelogin := false

	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			// got an operation error, so let's try to trigger the sso login!
			isRelogin = true
		} else {
			return fmt.Errorf("error while calling get-caller-identity with profile %s: %v", p.name, err)
		}
	}

	if isRelogin {
		if err := reloginWorkflow(&cfg, sharedConfig, p.configPath); err != nil {
			return err
		}
	} else {
		fmt.Printf("credentials for profile %s and role %s already available\n", p.name, aws.ToString(output.Arn))
	}

	return nil
}

func reloginWorkflow(cfg *aws.Config, sharedConfig config.SharedConfig, configPath string) error {
	// first let's build the cache filename
	cacheFileName, err := getCacheFileName(sharedConfig.SSOStartURL)

	if err != nil {
		return fmt.Errorf("error while generating cache filename: %v", err)
	}

	configBasePath := filepath.Dir(configPath)
	cacheFilePath := path.Join(configBasePath, ssoCachePath, cacheFileName)

	ssooidcClient := ssooidc.NewFromConfig(*cfg)

	register, err := ssooidcClient.RegisterClient(context.TODO(), &ssooidc.RegisterClientInput{
		ClientName: aws.String("sso-client-name"),
		ClientType: aws.String("public"),
		Scopes:     []string{"sso-portal:*"},
	})

	if err != nil {
		return fmt.Errorf("error while registering ssooidc client for profile %s: %v", sharedConfig.Profile, err)
	}

	deviceAuth, err := ssooidcClient.StartDeviceAuthorization(context.TODO(), &ssooidc.StartDeviceAuthorizationInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		StartUrl:     &sharedConfig.SSOStartURL,
	})

	if err != nil {
		return fmt.Errorf("error while starting device authorization for profile %s: %v", sharedConfig.Profile, err)
	}

	url := aws.ToString(deviceAuth.VerificationUriComplete)

	fmt.Printf("if browser is not opened automatically, please open link:\n%v\n", url)

	if err := browser.OpenURL(url); err != nil {
		return fmt.Errorf("error while opening url %s in browser: %v", url, err)
	}

	fmt.Println("Press ENTER key once login is done...")

	_ = bufio.NewScanner(os.Stdin).Scan()

	token, err := ssooidcClient.CreateToken(context.TODO(), &ssooidc.CreateTokenInput{
		ClientId:     register.ClientId,
		ClientSecret: register.ClientSecret,
		DeviceCode:   deviceAuth.DeviceCode,
		GrantType:    aws.String("urn:ietf:params:oauth:grant-type:device_code"),
	})

	if err != nil {
		return fmt.Errorf("error while generating token for profile %s: %v", sharedConfig.Profile, err)
	}

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second).UTC()
	registrationExpiresAt := time.Unix(register.ClientSecretExpiresAt, 0).UTC()

	tf := tokenFields{
		StartURL:              sharedConfig.SSOStartURL,
		Region:                sharedConfig.Region,
		AccessToken:           aws.ToString(token.AccessToken),
		ExpiresAt:             expiresAt.Format(time.RFC3339),
		ClientID:              aws.ToString(register.ClientId),
		ClientSecret:          aws.ToString(register.ClientSecret),
		RegistrationExpiresAt: registrationExpiresAt.Format(time.RFC3339),
	}

	asBytes, err := json.Marshal(&tf)

	if err != nil {
		return fmt.Errorf("error while marshaling token: %v", err)
	}

	if err := os.WriteFile(cacheFilePath, asBytes, 0600); err != nil {
		return fmt.Errorf("error while saving sso cache credentials file %s: %v", cacheFilePath, err)
	}

	return nil
}

func getCacheFileName(url string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(url))
	if err != nil {
		return "", err
	}
	return strings.ToLower(hex.EncodeToString(hash.Sum(nil))) + ".json", nil
}

func validateSharedConfig(sharedConfig config.SharedConfig) error {
	if sharedConfig.SSOStartURL == "" {
		return fmt.Errorf("error, sso_start_url field not found for profile %s", sharedConfig.Profile)
	}

	if sharedConfig.SSOAccountID == "" {
		return fmt.Errorf("error, sso_account_id field not found for profile %s", sharedConfig.Profile)
	}

	if sharedConfig.SSORoleName == "" {
		return fmt.Errorf("error, sso_role_name field not found for profile %s", sharedConfig.Profile)
	}

	return nil
}

func getSharedConfig(cfg *aws.Config) config.SharedConfig {
	var sharedConfig config.SharedConfig

	cfgSources := cfg.ConfigSources

	for _, cfgSource := range cfgSources {
		if foundSharedConfig, ok := cfgSource.(config.SharedConfig); ok {
			sharedConfig = foundSharedConfig
			break
		}
	}

	return sharedConfig
}
