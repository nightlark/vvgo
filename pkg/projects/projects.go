package projects

import (
	"errors"
	"github.com/virtual-vgo/vvgo/pkg/log"
)

var logger = log.Logger()

var ErrNotFound = errors.New("project not found")

// A VVGO project
type Project struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Released     bool     `json:"released"`
	Archived     bool     `json:"archived"`
	Sources      []string `json:"sources"`
	Composers    []string `json:"composers"`
	Arrangers    []string `json:"arrangers"`
	Editors      []string `json:"editors"`
	Transcribers []string `json:"transcribers"`
	Preparers    []string `json:"preparers"`
	ClixBy       []string `json:"clix_by"`
	Reviewers    []string `json:"reviewers"`
	Lyricists    []string `json:"lyricists"`
	AddlContent  []string `json:"addl_content"`
}

var project = Projects{projects: []Project{
	{
		Name:         "01-snake-eater",
		Title:        "Snake Eater",
		Released:     true,
		Archived:     false,
		Sources:      []string{"Metal Gear Solid 3"},
		Composers:    []string{"Norihiko Hibino"},
		Arrangers:    []string{},
		Editors:      []string{"Jerome Landingin"},
		Transcribers: []string{},
		Preparers:    []string{"The Giggling Donkey, Inc."},
		ClixBy:       []string{"Finny Jacob Zeleny"},
		Reviewers:    []string{},
		Lyricists:    []string{},
		AddlContent:  []string{"B. Harnish"},
	},
	{
		Name:         "02-proof-of-a-hero",
		Title:        "Proof of a Hero",
		Released:     true,
		Archived:     false,
		Sources:      []string{"Monster Hunter"},
		Composers:    []string{},
		Arrangers:    []string{"Jacob Zeleny"},
		Editors:      []string{},
		Transcribers: []string{},
		Preparers:    []string{"The Giggling Donkey, Inc.", "Thomas Håkanson"},
		ClixBy:       []string{"Finny Jacob Zeleny"},
		Reviewers:    []string{"B. Harnish"},
		Lyricists:    []string{"B. Harnish"},
		AddlContent:  []string{"Chris Suzuki", "B. Harnish"},
	},
}}

func init() {
	// build indices
	project.nameIndex = make(Index)
	for i, p := range project.projects {
		project.nameIndex[p.Name] = &project.projects[i]
	}
}

func GetName(name string) *Project { return project.GetName(name) }
func Exists(name string) bool      { return project.Exists(name) }

type Projects struct {
	projects  []Project
	nameIndex Index
}

type Index map[string]*Project

func (x *Projects) GetName(name string) *Project {
	return x.nameIndex[name]
}

func (x *Projects) Exists(name string) bool {
	_, ok := x.nameIndex[name]
	return ok
}
