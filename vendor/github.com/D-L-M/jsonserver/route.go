package jsonserver

import (
	"strings"
)

// Route structs define executable HTTP routes
type Route struct {
	Path       string
	Action     RouteAction
	Middleware []Middleware
}

// MatchesPath checks whether the route's path matches a given path and returns any wildcard values
func (route *Route) MatchesPath(path string) (bool, RouteParams) {

	normalisedRoutePath := normalisePath(route.Path)
	normalisedPath := normalisePath(path)
	pathFragments := strings.Split(normalisedPath, "/")
	routePathFragments := strings.Split(normalisedRoutePath, "/")
	hasFinalWildcard := strings.HasSuffix(normalisedRoutePath, "/:") || normalisedRoutePath == ":"
	lengthMatches := len(pathFragments) == len(routePathFragments)
	lengthMatchesWithFinalWildcard := hasFinalWildcard && len(pathFragments) >= len(routePathFragments)
	wildcardValues := RouteParams{}

	if lengthMatches || lengthMatchesWithFinalWildcard {

		for i, routePathFragment := range routePathFragments {

			isWildcard := strings.HasPrefix(routePathFragment, "{") && strings.HasSuffix(routePathFragment, "}")
			isFinalWildcard := hasFinalWildcard && i == (len(routePathFragments)-1)

			// The route path no longer matches
			if isWildcard == false && isFinalWildcard == false && pathFragments[i] != routePathFragment {
				return false, RouteParams{}
			}

			// The route matches on a final wildcard, so compile the remaining route param values
			if isFinalWildcard {

				wildcardValues["{catchAll}"] = strings.Join(pathFragments[i:], "/")

				// The route matches on a wildcard, so obtain its key and value
			} else if isWildcard {

				wildcardKey := routePathFragment[1 : len(routePathFragment)-1]
				wildcardValues[wildcardKey] = pathFragments[i]

			}

		}

		return true, wildcardValues

	}

	return false, RouteParams{}

}

// normalisePath ensures that a path has no leading or trailing slashes
func normalisePath(path string) string {

	if strings.HasPrefix(path, "/") == false {
		path = "/" + path
	}

	if strings.HasSuffix(path, "/") == false {
		path += "/"
	}

	if path == "/" {
		return ""
	}

	return path[1 : len(path)-1]

}
