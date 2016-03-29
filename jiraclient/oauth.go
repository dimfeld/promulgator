package jiraclient

// Note: This is currently not used, in favor of HTTP Basic Auth.
//
// import (
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"net/http"
//
// 	"github.com/mrjones/oauth"
// )
//
//
// // NewOAuthClient returns an HTTP client that is configured to use the
// // given access token and secret. If you do not have an access token, use
// // a utility such as curlicue to perform the OAuth dance and get one.
// func NewOAuthClient(key, accessToken, accessSecret, URLBase string) (*http.Client, error) {
// 	// These URLs are JIRA-specific
// 	requestURL := URLBase + "plugins/servlet/oauth/request-token"
// 	authorizeURL := URLBase + "plugins/servlet/oauth/authorize"
// 	accessURL := URLBase + "plugins/servlet/oauth/access-token"
//
// 	sp := oauth.ServiceProvider{
// 		RequestTokenUrl:   requestURL,
// 		AuthorizeTokenUrl: authorizeURL,
// 		AccessTokenUrl:    accessURL,
// 		HttpMethod:        "POST",
// 	}
//
// 	// Just generate a new key every time. We don't need it to be the same
// 	// between runs.
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 256)
// 	if err != nil {
// 		return nil, err
// 	}
// 	consumer := oauth.NewRSAConsumer(key, privateKey, sp)
// 	tokenStruct := &oauth.AccessToken{
// 		Token:  accessToken,
// 		Secret: accessSecret,
// 	}
// 	return consumer.MakeHttpClient(tokenStruct)
// }
