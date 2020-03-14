package loader

import (
	// "fmt"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/xerrors"
	"regexp"
)

var InvalidVersionValue = xerrors.New("Invalid version value")

type _version string

func (v _version) MINOR() int {
	if err := v.Verify(); err != nil {
		panic(err)
	}
	if val, err := strconv.Atoi(strings.Split(string(v), ".")[0]); err != nil {
		panic(xerrors.Errorf("Extract Minor raise error %s:%w", v, err))
	} else {
		return val
	}
}

func (v _version) PATCH() int {
	if err := v.Verify(); err != nil {
		panic(err)
	}
	if val, err := strconv.Atoi(strings.Split(string(v), ".")[1]); err != nil {
		panic(xerrors.Errorf("Extract PATCH raise error %s:%w", v, err))
	} else {
		return val
	}
}

var versionReg, _ = regexp.Compile("^[0-9]+\\.[0-9]+$")

func (v _version) Verify() error {
	if !versionReg.Match([]byte(v)) {
		return xerrors.Errorf("Raise error of value \"%s\":%w", v, InvalidVersionValue)
	}
	return nil

}

type Version struct {
	MINOR int
	PATCH int
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.MINOR, v.PATCH)
}

func NewVersion(value string) (Version, error) {
	_v := _version(value)
	if err := _v.Verify(); err != nil {
		return Version{}, err
	}
	return Version{_v.MINOR(), _v.PATCH()}, nil
}

func NewVersionPanic(value string) Version {
	_v := _version(value)
	if err := _v.Verify(); err != nil {
		panic(err)
	}
	return Version{_v.MINOR(), _v.PATCH()}
}

func (v Version) MarshalBSONValue() (bsontype.Type, []byte, error) {

	return bson.MarshalValue(v.String())
}

func (v *Version) UnmarshalBSONValue(t bsontype.Type, src []byte) error {
	// s := ""
	var err error
	var val bsonx.Val
	// val.UnmarshalBSONValue(t,src)
	if err := val.UnmarshalBSONValue(t, src); err != nil {
		return err
	}

	if *v, err = NewVersion(val.String()); err != nil {
		return err
	}
	return nil
}

func (v Version) NextMINOR() Version {
	return Version{v.MINOR + 1, v.PATCH}
}
func (v Version) NextPATCH() Version {
	return Version{v.MINOR, v.PATCH + 1}
}

func (v Version) Less(anotherVersion Version) bool {
	if v.MINOR < anotherVersion.MINOR {
		return true
	} else if v.MINOR == anotherVersion.MINOR && v.PATCH < anotherVersion.PATCH {
		return true
	}
	return false
}

func (v Version) Greater(anotherVersion Version) bool {
	if v.MINOR > anotherVersion.MINOR {
		return true
	} else if v.MINOR == anotherVersion.MINOR && v.PATCH > anotherVersion.PATCH {
		return true
	}
	return false
}

type Versions []Version

func (vs *Versions) Max() *Version {
	var max Version
	if cap(*vs) == 0 {
		return nil
	}
	for _, v := range *vs {
		if v.Greater(max) {
			max = v
		}
	}
	return &max
}

type HasVersion interface {
	GetVersion() Version
	SetVersion(Version)
	GetData() interface{}
	SetData(interface{})
}