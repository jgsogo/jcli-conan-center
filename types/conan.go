// Package types provides basic types to work with Conan items: references and packages.
package types

import (
	"fmt"
	"strings"
	"regexp"
)

// Constants to be used with Conan elements.
const (
	ValidConanChars       = `[a-zA-Z0-9_][a-zA-Z0-9_\+\.-]`  // Validates (regex) a part from a Conan reference
	FilesystemPlaceHolder = "_"  // Filesystem representation of a null user or null channel in a Conan reference. 
)

// ParseStringReference parses a string and returns a Conan reference. The string can have any of these formats: name/version,
// name/version@user/channel, name/version#revision, name/version@user/channel#revision.
func ParseStringReference(reference string) (*Reference, error) {
	referencePattern := regexp.MustCompile(`^(?P<name>` + ValidConanChars + `+)\/(?P<version>` + ValidConanChars + `+)(@(?P<user>` + ValidConanChars + `+)\/(?P<channel>` + ValidConanChars + `+))?(#(?P<revision>[a-z0-9]+))?$`)
	m := referencePattern.FindStringSubmatch(reference)
	if m == nil {
		return nil, fmt.Errorf("String '%s' doesn't match a Conan reference", reference)
	}
	name := m[1]
	version := m[2]
	user := m[4]
	channel := m[5]
	revision := m[7]

	if user == "" || channel == "" {
		return &Reference{Name: name, Version: version, User: nil, Channel: nil, Revision: revision}, nil
	}
	return &Reference{Name: name, Version: version, User: &user, Channel: &channel, Revision: revision}, nil
}

// Reference represents a Conan reference with its parts: name, version, user, channel and revision. Only the attributes
// `Channel` and `User` are optional in a valid reference.
type Reference struct {
	Name     string
	Version  string
	User     *string
	Channel  *string
	Revision string
}

// ToString returns a string representation of the `Reference`. Use the argument `withRevision` to add or not the
// revision to the output.
func (ref *Reference) ToString(withRevision bool) string {
	var ret string
	if ref.User != nil {
		ret = fmt.Sprintf("%s/%s@%s/%s", ref.Name, ref.Version, *ref.User, *ref.Channel)
	} else {
		ret = fmt.Sprintf("%s/%s", ref.Name, ref.Version)
	}
	if withRevision {
		ret = ret + "#" + ref.Revision
	}
	return ret
}

// String converts to string the reference (with revisions).
func (ref *Reference) String() string {
	return ref.ToString(true)
}

// RtPath returns the path inside Artifactory to the `Reference`. It can be considered with
// or without revisions (latest element in the Artifactory path). 
func (ref *Reference) RtPath(withRevision bool) string {
	var user string
	if ref.User == nil {
		user = FilesystemPlaceHolder
	} else {
		user = *ref.User
	}

	var channel string
	if ref.Channel == nil {
		channel = FilesystemPlaceHolder
	} else {
		channel = *ref.Channel
	}

	str := []string{user, ref.Name, ref.Version, channel}
	if withRevision {
		str = append(str, ref.Revision)
	}
	return strings.Join(str, "/")
}

// Package represents a Conan package with its `Reference`, the package ID and the package revision.
type Package struct {
	Ref       Reference
	PackageId string
	Revision  string
}

func (pkg *Package) String() string {
	return pkg.ToString(true)
}

// ToString returns the string representation of the Conan package. Use the argument `withRevision` to add or not the
// revisions to the output.
func (pkg *Package) ToString(withRevision bool) string {
	ret := pkg.Ref.ToString(withRevision) + ":" + pkg.PackageId
	if withRevision {
		ret = ret + "#" + pkg.Revision
	}
	return ret
}

// RtPath returns the path inside Artifactory to the `Package`. It can be considered with
// or without revisions (last element in the Artifactory path). However, not that the recipe
// will always contain the revision, as it is an element in the middle of the path. 
func (pkg *Package) RtPath(withRevision bool) string {
	str := []string{pkg.Ref.RtPath(true), "package", pkg.PackageId}
	if withRevision {
		str = append(str, pkg.Revision)
	}
	return strings.Join(str, "/")
}
