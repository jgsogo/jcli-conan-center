package types

import (
	"fmt"
	"strings"
	"regexp"
)

const (
	ValidConanChars       = `[a-zA-Z0-9_][a-zA-Z0-9_\+\.-]`
	FilesystemPlaceHolder = "_"
)

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
		if channel != "" || user != "" {
			panic("Provided reference contains 'channel' or 'user', but not both!")
		}
		return &Reference{Name: name, Version: version, User: nil, Channel: nil, Revision: revision}, nil
	}
	return &Reference{Name: name, Version: version, User: &user, Channel: &channel, Revision: revision}, nil
}

type Reference struct {
	Name     string
	Version  string
	User     *string
	Channel  *string
	Revision string
}

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

func (ref *Reference) String() string {
	return ref.ToString(true)
}

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

type Package struct {
	Ref       Reference
	PackageId string
	Revision  string
}

func (pkg *Package) String() string {
	return pkg.ToString(true)
}

func (pkg *Package) ToString(withRevision bool) string {
	ret := pkg.Ref.ToString(withRevision) + ":" + pkg.PackageId
	if withRevision {
		ret = ret + "#" + pkg.Revision
	}
	return ret
}

func (pkg *Package) RtPath(withRevision bool) string {
	str := []string{pkg.Ref.RtPath(true), "package", pkg.PackageId}
	if withRevision {
		str = append(str, pkg.Revision)
	}
	return strings.Join(str, "/")
}
