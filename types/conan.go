package types

import (
	"fmt"
	"strings"
)

const (
	ValidConanChars       = `[a-zA-Z0-9_][a-zA-Z0-9_\+\.-]`
	FilesystemPlaceHolder = "_"
)

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
