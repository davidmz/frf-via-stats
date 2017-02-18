package main

import (
	"archive/zip"
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/juju/errors"
)

var (
	feedInfoRe = regexp.MustCompile(`^[a-z0-9-]+/_json/data/feedinfo\.js$`)
	entryRe    = regexp.MustCompile(`^[a-z0-9-]+/_json/data/entries/[0-9a-f]{8}\.js$`)
	twitterRe  = regexp.MustCompile(`^(http://twitter\.com/\w+)/statuses`)
)

func processArchive(archFile string) (*statInfo, error) {
	arch, err := zip.OpenReader(archFile)
	if err != nil {
		return nil, errors.Annotate(err, "cannot open archive")
	}
	defer arch.Close()

	si := new(statInfo)

	for _, f := range arch.File {
		if feedInfoRe.MatchString(f.Name) {
			username, err := processFeedInfo(f)
			if err != nil {
				return nil, errors.Annotate(err, "error in feedinfo")
			}
			si.UserName = username
			break
		}
	}

	if si.UserName == "" {
		return nil, errors.New("cannot find feedinfo")
	}

	for _, f := range arch.File {
		if entryRe.MatchString(f.Name) {
			src, err := processEntry(f)
			if err != nil {
				return nil, errors.Annotate(err, "cannot process entry")
			}
			found := false
			for _, s := range si.Sources {
				if s.URL == src.URL {
					found = true
					s.Count++
					break
				}
			}
			if !found {
				src.Count = 1
				si.Sources = append(si.Sources, src)
			}
		}
	}

	return si, nil
}

func processFeedInfo(file *zip.File) (string, error) {
	user := new(userJSON)
	if err := readZipEntry(file, user); err != nil {
		return "", err
	}

	if user.Type != "user" {
		return "", errors.Errorf("@%s is not a user (%s)", user.UserName, user.Type)
	}
	return user.UserName, nil
}

func processEntry(file *zip.File) (*srcInfo, error) {
	entry := new(entryJSON)
	if err := readZipEntry(file, entry); err != nil {
		return nil, err
	}

	if twitterRe.MatchString(entry.Via.URL) {
		// twitter is a special case
		entry.Via.URL = twitterRe.FindStringSubmatch(entry.Via.URL)[1]
	}

	si := new(srcInfo)
	si.Name = entry.Via.Name
	si.URL = entry.Via.URL

	if si.URL == "" {
		si.URL = "http://friendfeed.com"
		si.Name = "FriendFeed"
	}

	return si, nil
}

func readZipEntry(file *zip.File, v interface{}) error {
	r, err := file.Open()
	if err != nil {
		return errors.Annotate(err, "cannot open archived file")
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Annotate(err, "cannot read archived file")
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return errors.Annotate(err, "cannot parse entry")
	}

	return nil
}
