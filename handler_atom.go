package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/kjk/blog/pkg/notes"
	atom "github.com/thomas11/atomgenerator"
)

func copyAndSortArticles(articles []*Article) []*Article {
	n := len(articles)
	res := make([]*Article, n, n)
	copy(res, articles)
	sort.Slice(res, func(i, j int) bool {
		return res[j].PublishedOn.After(res[i].PublishedOn)
	})
	return res
}

func handleAtomHelp(w http.ResponseWriter, r *http.Request, excludeNotes bool) {
	articles := store.GetArticles(false)
	if excludeNotes {
		articles = filterArticlesByTag(articles, "note", false)
	}
	articles = copyAndSortArticles(articles)
	n := 25
	if n > len(articles) {
		n = len(articles)
	}

	latest := make([]*Article, n, n)
	size := len(articles)
	for i := 0; i < n; i++ {
		latest[i] = articles[size-1-i]
	}

	pubTime := time.Now()
	if len(articles) > 0 {
		pubTime = articles[0].PublishedOn
	}

	feed := &atom.Feed{
		Title:   "Krzysztof Kowalczyk blog",
		Link:    "https://blog.kowalczyk.info/atom.xml",
		PubDate: pubTime,
	}

	for _, a := range latest {
		//id := fmt.Sprintf("tag:blog.kowalczyk.info,1999:%d", a.Id)
		e := &atom.Entry{
			Title:   a.Title,
			Link:    "https://blog.kowalczyk.info" + a.URL(),
			Content: a.BodyHTML,
			PubDate: a.PublishedOn,
		}
		feed.AddEntry(e)
	}

	s, err := feed.GenXml()
	if err != nil {
		s = []byte("Failed to generate XML feed")
	}

	w.Write(s)
}

// /atom-all.xml
func handleAtomAll(w http.ResponseWriter, r *http.Request) {
	handleAtomHelp(w, r, false)
}

// /atom.xml
func handleAtom(w http.ResponseWriter, r *http.Request) {
	handleAtomHelp(w, r, true)
}

// /dailynotes-atom.xml
// TODO: could cache generated xml
func handleNotesFeed(w http.ResponseWriter, r *http.Request) {
	notes := notes.NotesAllNotes
	if len(notes) > 25 {
		notes = notes[:25]
	}

	pubTime := time.Now()
	if len(notes) > 0 {
		pubTime = notes[0].Day
	}

	feed := &atom.Feed{
		Title:   "Krzysztof Kowalczyk daily notes",
		Link:    "https://blog.kowalczyk.info/dailynotes-atom.xml",
		PubDate: pubTime,
	}

	for _, n := range notes {
		//id := fmt.Sprintf("tag:blog.kowalczyk.info,1999:%d", a.Id)
		title := n.Title
		if title == "" {
			title = n.ID
		}
		e := &atom.Entry{
			Title:   title,
			Link:    "https://blog.kowalczyk.info/" + n.URL,
			Content: string(n.HTMLBody),
			PubDate: n.Day,
		}
		feed.AddEntry(e)
	}

	s, err := feed.GenXml()
	if err != nil {
		s = []byte("Failed to generate XML feed")
	}

	w.Write(s)
}
