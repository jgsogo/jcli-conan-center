package types

import (
	"fmt"
	"strings"
)

const (
	ValidConanChars = `[a-zA-Z0-9_][a-zA-Z0-9_\+\.-]`
)

type Reference struct {
	Name     string
	Version  string
	User     string
	Channel  string
	Revision string
}

func (ref *Reference) ToString(withRevision bool) string {
	var ret string
	if len(ref.User) > 0 {
		ret = fmt.Sprintf("%s/%s@%s/%s", ref.Name, ref.Version, ref.User, ref.Channel)
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
	user := ref.User
	if len(user) == 0 {
		user = "_"
	}
	channel := ref.Channel
	if len(channel) == 0 {
		channel = "_"
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
	return fmt.Sprintf("%s:%s#%s", pkg.Ref, pkg.PackageId, pkg.Revision)
}

//func (pkg *Package) RtPath() string {
//	str := []string{pkg.Ref.rtPath(), "package", pkg.PackageId, pkg.Revision}
//	return strings.Join(str, "/")
//}
