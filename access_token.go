package wego

import (
	"crypto/md5"
	"fmt"
	"github.com/godcong/wego/cache"
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"strings"
	"time"
)

// AccessTokenOption ...
type AccessTokenOption struct {
	RemoteHost string
	TokenKey   string
	TokenURL   string
}

/*AccessToken AccessToken */
type AccessToken struct {
	*AccessTokenConfig
	remoteHost string
	tokenKey   string
	tokenURL   string
}

/*AccessTokenSafeSeconds token安全时间 */
const AccessTokenSafeSeconds = 500

// RemoteHost ...
func (obj *AccessToken) RemoteHost() string {
	if obj != nil && obj.remoteHost != "" {
		return obj.remoteHost
	}
	return apiWeixin
}

// TokenURL ...
func (obj *AccessToken) TokenURL() string {
	return util.URL(obj.RemoteHost(), tokenURL(obj))
}
func tokenURL(obj *AccessToken) string {
	if obj != nil && obj.tokenURL != "" {
		return obj.tokenURL
	}
	return accessTokenURLSuffix
}

/*NewAccessToken NewAccessToken*/
func NewAccessToken(property *AccessTokenConfig, opts ...*AccessTokenOption) *AccessToken {
	token := &AccessToken{
		AccessTokenConfig: property,
	}
	token.parse(opts)
	return token
}

func (obj *AccessToken) parse(opts []*AccessTokenOption) {
	if opts == nil {
		return
	}
	obj.remoteHost = opts[0].RemoteHost
	obj.tokenURL = opts[0].TokenURL
	obj.tokenKey = opts[0].TokenKey
}

/*Refresh 刷新AccessToken */
func (obj *AccessToken) Refresh() *AccessToken {
	log.Debug("AccessToken|Refresh")
	obj.getToken(true)
	return obj
}

/*GetRefreshToken 获取刷新token */
func (obj *AccessToken) GetRefreshToken() *core.Token {
	log.Debug("AccessToken|GetRefreshedToken")
	return obj.getToken(true)
}

/*GetToken 获取token */
func (obj *AccessToken) GetToken() *core.Token {
	return obj.getToken(false)
}

// KeyMap ...
func (obj *AccessToken) KeyMap() util.Map {
	return MustKeyMap(obj)
}

func (obj *AccessToken) getToken(refresh bool) *core.Token {
	key := obj.getCacheKey()

	if !refresh && cache.Has(key) {
		if v, b := cache.Get(key).(*core.Token); b {
			if v.ExpiresIn > time.Now().Unix() {
				log.Debugf("cached accessToken found:%+v", v)
				return v
			}
		}
	}

	token, e := requestToken(obj.TokenURL(), obj.AccessTokenConfig)
	if e != nil {
		log.Error(e)
		return nil
	}

	log.Debug("getToken:%+v", *token)
	if v := token.ExpiresIn; v != 0 {
		obj.SetTokenWithLife(token.AccessToken, v-AccessTokenSafeSeconds)
	} else {
		obj.SetToken(token.AccessToken)
	}
	return token
}

func requestToken(url string, credentials *AccessTokenConfig) (*core.Token, error) {
	var token core.Token
	e := Get(url, credentials.ToMap()).Unmarshal(&token)
	if e != nil {
		return nil, e
	}
	return &token, nil
}

/*SetTokenWithLife set string accessToken with life time */
func (obj *AccessToken) SetTokenWithLife(token string, lifeTime int64) *AccessToken {
	return obj.setToken(token, lifeTime)
}

/*SetToken set string accessToken */
func (obj *AccessToken) SetToken(token string) *AccessToken {
	return obj.setToken(token, 7200-AccessTokenSafeSeconds)
}

func (obj *AccessToken) setToken(token string, lifeTime int64) *AccessToken {
	cache.SetWithTTL(obj.getCacheKey(), &core.Token{
		AccessToken: token,
		ExpiresIn:   time.Now().Add(time.Duration(lifeTime)).Unix(),
	}, lifeTime)
	return obj
}

func (obj *AccessToken) getCredentials() string {
	cred := strings.Join([]string{obj.GrantType, obj.AppID, obj.AppSecret}, ".")
	c := md5.Sum([]byte(cred))
	return fmt.Sprintf("%x", c[:])
}

func (obj *AccessToken) getCacheKey() string {
	return "godcong.wego.access_token." + obj.getCredentials()
}

const accessTokenNil = "nil point accessToken"
const tokenNil = "nil point token"

/*MustKeyMap get accessToken's key,value with map when nil or error return nil map */
func MustKeyMap(at *AccessToken) util.Map {
	m := util.Map{}
	if m, e := KeyMap(at); e != nil {
		return m
	}
	return m
}

/*KeyMap get accessToken's key,value with map */
func KeyMap(at *AccessToken) (util.Map, error) {
	if at == nil {
		return nil, xerrors.New(accessTokenNil)
	}
	if token := at.GetToken(); token != nil {
		return token.KeyMap(), nil
	}
	return nil, xerrors.New(tokenNil)
}
