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

	assert.Equal(t, conanPackage.String(), "name/version#rrev:pkgId#prev")
}
