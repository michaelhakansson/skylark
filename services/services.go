package services

import (
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/structures"
)

// Service describes the interface for a service
type Service interface {
    GetAllProgramIDs() []string
    GetShow(string) (structures.Show, []string)
    GetEpisode(string) structures.Episode
    GetName() string
}

var services []Service

func init() {
    svt := svtplay.SVTPlay{}
    tv3 := tv3play.TV3Play{}
    kanal5 := kanal5play.Kanal5Play{}
    services = []Service{svt, tv3, kanal5}
}

// GetServices returns a list of the supported services
func GetServices() []string {
    var servicelist []string
    for _, service := range services {
        servicelist = append(servicelist, service.GetName())
    }
    return servicelist
}

// GetIdsWithService returns the result of the 'GetAllProgramIds' function
// of the specified service
func GetIdsWithService(playservice string) (IDs []string) {
    for _, service := range services {
        if service.GetName() == playservice {
            IDs = service.GetAllProgramIDs()
            return
        }
    }
    return
}

// GetShowWithService returns the result of the 'GetShow' function of the
// specified service
func GetShowWithService(showid string, playservice string) (show structures.Show, episodeIDs []string) {
    for _, service := range services {
        if service.GetName() == playservice {
            show, episodeIDs = service.GetShow(showid)
            return
        }
    }
    return
}

// GetEpisodeWithService returns the result of the 'GetEpisode' function of the
// specified service
func GetEpisodeWithService(episodeid string, playservice string) (episode structures.Episode) {
    for _, service := range services {
        if service.GetName() == playservice {
            episode = service.GetEpisode(episodeid)
            return
        }
    }
    return
}
