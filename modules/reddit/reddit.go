package reddit

import (
	"github.com/jzelinskie/geddit"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/logging"
	"github.com/pkg/errors"
	"math/rand"
	"strings"
	"time"
)

var log = logging.New("reddit")

const pollInterval = 15

var lastLink = ""

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	o, err := redditAuth(conf)
	if err != nil {
		return nil, err
	}
	rv := &reddit{
		Base: core.NewBase("reddit", sender, conf),
		o:    o,
	}
	go rv.startPolling(conf, o)
	rv.AddCommand("le", rv.le)
	rv.AddCommand("lepic", rv.lepic)
	//	rv.AddCommand("lelong", rv.lelong)
	rv.AddCommand("lelast", rv.lelast)
	return rv, nil
}

func redditAuth(conf config.Config) (*geddit.OAuthSession, error) {
	log.Info("Starting new Reddit OAuth session...")
	o, err := geddit.NewOAuthSession(
		conf.Reddit.ClientID,
		conf.Reddit.ClientSecret,
		"Testing OAuth Bot by u/my_user v0.1 see source https://github.com/jzelinskie/geddit",
		"http://redirect.url",
	)
	if err != nil {
		return nil, err
	}
	// Create new auth token for confidential clients (personal scripts/apps).
	log.Info("Auth succeeded, logging in...")
	err = o.LoginAuth(conf.Reddit.Username, conf.Reddit.Password)
	if err != nil {
		return nil, err
	}
	log.Info("Done logging into reddit...")
	return o, nil
}

func (r *reddit) startPolling(conf config.Config, o *geddit.OAuthSession) (*geddit.OAuthSession, error) {
	tokenTime := time.Now()
	var err error
	for {
		if time.Since(tokenTime).Hours() >= 1 {
			log.Info("Refreshing token...")
			r.o, err = redditAuth(conf)
			if err != nil {
				return nil, err
			}
			tokenTime = time.Now()
		}
		time.Sleep(pollInterval * time.Second)
	}
}

type reddit struct {
	core.Base
	o           *geddit.OAuthSession
	commentList []*geddit.Comment
}

func (r *reddit) getSubOrSelectSub(arguments core.CommandArguments) (string, error) {
	if len(arguments.Arguments) > 0 {
		return arguments.Arguments[0], nil
	} else {
		return "linuxcirclejerk", nil
	}
}

func (r *reddit) le(arguments core.CommandArguments) ([]string, error) {
	err := r.getSubComments(arguments)
	if err != nil {
		return nil, err
	}
	var comm *geddit.Comment
	comm, err = r.getRandomComment(r.commentList)
	if err != nil {
		return nil, err
	}
	lastLink = comm.Permalink
	if len(comm.Body)>99{
		return []string{comm.Body[:99]}, nil
	}
	return []string{comm.Body}, nil
}

func (r *reddit) lelast(arguments core.CommandArguments) ([]string, error) {
	return []string{"https://reddit.com" + lastLink}, nil
}

//func (r *reddit) lelong(arguments core.CommandArguments) ([]string, error) {
//	err := r.getSubComments(arguments)
//	if err != nil {
//		return nil, err
//	}
//	currcomm, err := r.getRandomComment(r.commentList)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Printf("length of str %d\n", len(strings.Join(currcomm, " ")))
//	if len(strings.Join(currcomm, " ")) > 140 {
//		return strings.Split(strings.Join(currcomm, " "), "\n"), nil
//	} else {
//		r.lelong(arguments)
//	}
//	return nil, errors.New("Could not find a random post")
//}

func (r *reddit) getSubComments(arguments core.CommandArguments) error {
	sub, _ := r.getSubOrSelectSub(arguments)
	var err error
	r.commentList, err = r.o.SubredditComments(sub)
	if err != nil {
		return err
	}
	return nil
}

func (r *reddit) getRandomComment(commentlist []*geddit.Comment) (*geddit.Comment, error) {
	if len(commentlist) > 0 {
		commentIndex := rand.Intn(len(commentlist))
		return commentlist[commentIndex], nil
	}
	return nil, errors.New("Could not find a random post")
}

func (r *reddit) lepic(arguments core.CommandArguments) ([]string, error) {
	sub, _ := r.getSubOrSelectSub(arguments)
	subOpts := geddit.ListingOptions{
		Limit: 500,
	}
	posts, err := r.o.SubredditSubmissions(sub, geddit.DefaultPopularity, subOpts)
	if err != nil {
		return nil, err
	}
	if len(posts) > 0 {
		commentIndex := rand.Intn(len(posts))
		if len(posts[commentIndex].URL) > 0 {
			return strings.Split(posts[commentIndex].URL, "\n"), nil
		} else {
			r.lepic(arguments)
		}
	}
	return nil, errors.New("Could not find a random post")
}
