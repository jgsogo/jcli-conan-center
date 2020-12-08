package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReference(t *testing.T) {
	reference := Reference{"name", "version", nil, nil, "rrev"}
	assert.Equal(t, reference.Name, "name")
	assert.Equal(t, reference.Version, "version")
	assert.Nil(t, reference.User)
	assert.Nil(t, reference.Channel)
	assert.Equal(t, reference.Revision, "rrev")

	assert.Equal(t, reference.ToString(true), "name/version#rrev")
	assert.Equal(t, reference.ToString(false), "name/version")
	assert.Equal(t, reference.String(), "name/version#rrev")
	assert.Equal(t, reference.RtPath(true), "_/name/version/_/rrev")
	assert.Equal(t, reference.RtPath(false), "_/name/version/_")
}
func TestReferenceUserChannel(t *testing.T) {
	user := "user"
	channel := "channel"
	reference := Reference{"name", "version", &user, &channel, "rrev"}
	assert.Equal(t, reference.Name, "name")
	assert.Equal(t, reference.Version, "version")
	assert.Equal(t, *reference.User, "user")
	assert.Equal(t, *reference.Channel, "channel")
	assert.Equal(t, reference.Revision, "rrev")

	assert.Equal(t, reference.ToString(true), "name/version@user/channel#rrev")
	assert.Equal(t, reference.ToString(false), "name/version@user/channel")
	assert.Equal(t, reference.String(), "name/version@user/channel#rrev")
	assert.Equal(t, reference.RtPath(true), "user/name/version/channel/rrev")
	assert.Equal(t, reference.RtPath(false), "user/name/version/channel")
}

func TestPackage(t *testing.T) {
	reference := Reference{"name", "version", nil, nil, "rrev"}
	conanPackage := Package{Ref: reference, PackageId: "pkgId", Revision: "prev"}

	assert.Equal(t, conanPackage.ToString(true), "name/version#rrev:pkgId#prev")
	assert.Equal(t, conanPackage.ToString(false), "name/version:pkgId")
	assert.Equal(t, conanPackage.String(), "name/version#rrev:pkgId#prev")
	assert.Equal(t, conanPackage.RtPath(true), "_/name/version/_/rrev/package/pkgId/prev")
	assert.Equal(t, conanPackage.RtPath(false), "_/name/version/_/rrev/package/pkgId")
}

func TestParseStringReference(t *testing.T) {
	ref, err := ParseStringReference("name/version")
	assert.Nil(t, err)
	assert.Equal(t, "name", ref.Name)
	assert.Equal(t, "version", ref.Version)
	assert.Nil(t, ref.User)
	assert.Nil(t, ref.Channel)
	assert.Equal(t, "", ref.Revision)

	ref, err = ParseStringReference("name/version@user/channel")
	assert.Nil(t, err)
	assert.Equal(t, "name", ref.Name)
	assert.Equal(t, "version", ref.Version)
	assert.Equal(t, "user", *ref.User)
	assert.Equal(t, "channel", *ref.Channel)
	assert.Equal(t, "", ref.Revision) // TODO: Think about a better value when it has no value.

	ref, err = ParseStringReference("name/version#rrev")
	assert.Nil(t, err)
	assert.Equal(t, "name", ref.Name)
	assert.Equal(t, "version", ref.Version)
	assert.Nil(t, ref.User)
	assert.Nil(t, ref.Channel)
	assert.Equal(t, "rrev", ref.Revision)

	ref, err = ParseStringReference("name/version@user/channel#rrev")
	assert.Nil(t, err)
	assert.Equal(t, "name", ref.Name)
	assert.Equal(t, "version", ref.Version)
	assert.Equal(t, "user", *ref.User)
	assert.Equal(t, "channel", *ref.Channel)
	assert.Equal(t, "rrev", ref.Revision)
}

func TestParseStringReferenceErrors(t *testing.T) {
	_, err := ParseStringReference("name/version@")
	assert.NotNil(t, err)
	assert.Equal(t, "String 'name/version@' doesn't match a Conan reference", err.Error())

	_, err = ParseStringReference("name/version@#rrev")
	assert.NotNil(t, err)
	assert.Equal(t, "String 'name/version@#rrev' doesn't match a Conan reference", err.Error())

	_, err = ParseStringReference("n/version@")
	assert.NotNil(t, err)
	assert.Equal(t, "String 'n/version@' doesn't match a Conan reference", err.Error())
}
