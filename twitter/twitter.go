package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	twicli "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/okoshiyoshinori/twigolf-server/config"
	"github.com/okoshiyoshinori/twigolf-server/db"
	"github.com/okoshiyoshinori/twigolf-server/model"
	"github.com/okoshiyoshinori/twigolf-server/util"
)

type Token struct {
  AccessToken string
  AccessSecretToken string
}

func NewClient(userToken Token) *twicli.Client  {
  ck := config.GetConfig().Apserver.CsKey
  csk := config.GetConfig().Apserver.CsSecretKey
  uconf := oauth1.NewConfig(ck,csk)
  token := oauth1.NewToken(userToken.AccessToken,userToken.AccessSecretToken)
  httpClient := uconf.Client(oauth1.NoContext,token)
  return twicli.NewClient(httpClient)
}

func NewDmEvent(id string,massage string) *twicli.DirectMessageEventsNewParams {
  return &twicli.DirectMessageEventsNewParams{
    Event: &twicli.DirectMessageEvent{
      Type: "message_create",
      Message: &twicli.DirectMessageEventMessage{
        Target: &twicli.DirectMessageTarget{
          RecipientID: id,
        },
        Data: &twicli.DirectMessageData{
          Text:massage,
        },
      },
    },
  }
}

func GetConnect() *oauth.Client {
  return &oauth.Client{
    TemporaryCredentialRequestURI: config.GetConfig().Apserver.TwitterTempCredentialRequestURI,
    ResourceOwnerAuthorizationURI: config.GetConfig().Apserver.TwitterResourceOwnerAuthorizationURI,
    TokenRequestURI: config.GetConfig().Apserver.TwitterTokenRequestURI,
    Credentials: oauth.Credentials{
      Token: config.GetConfig().Apserver.CsKey,
      Secret: config.GetConfig().Apserver.CsSecretKey,
    },
  }
}

func GetAccessToken(rt *oauth.Credentials, oauthVerifer string) (*oauth.Credentials, error) {
	oc := GetConnect()
	at, _, err := oc.RequestToken(nil, rt, oauthVerifer)
	return at, err
}

func GetAccount(at *oauth.Credentials, user *model.TwitterAccount) error {
	oc := GetConnect()

	v := url.Values{}
	v.Set("include_email", "false")

	resp, err := oc.Get(nil, at,config.GetConfig().Apserver.TwitterVerifyCredentialsURI, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return errors.New("Twitter is unavailable")
	}

	if resp.StatusCode >= 400 {
		return errors.New("Twitter request is invalid")
	}

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		return err
	}
	return nil
}

func MakeOgpInfo(cid int) *model.Ogp {
  d := db.GetConn()
  competition := model.Competition{}
  if d.Where("id = ?",cid).Preload("User").Find(&competition).RecordNotFound() {
    return nil
  }
  //title
  title := competition.Title
  //type 
  typeStr := "article"
  //url
  url := fmt.Sprintf("%s/events/%d",config.GetConfig().Apserver.Host,competition.ID)

  var place string = ""
  if competition.PlaceText != nil {
    place = fmt.Sprintf("開催場所:%s",*competition.PlaceText)
  } else {
    place = "開催場所:未定"
  }

  var event_day string = ""
  if competition.EventDay != nil {
    event_day = fmt.Sprintf("開催日時:%s",util.DateToString(competition.EventDay))
  } else {
    event_day = "開催日時:未定"
  }
  //description
  description := place + "\n" + event_day

  twitter := model.TwitterCard{
    Card: "summary",
    Site: "@" + config.GetConfig().Apserver.TwitterAccount,
    Creater: "@" + competition.User.ScreenName,
  }

  ogp := model.Ogp{
    OgTitle: title,
    OgType: typeStr,
    OgDescription: description,
    OgSiteName: config.GetConfig().Apserver.SiteName,
    OgUrl: url,
    OgImage: config.GetConfig().Apserver.Image,
    Twitter: &twitter,
  }
  return &ogp
}
