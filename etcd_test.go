package etcd

import (
	"code.google.com/p/go-uuid/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEtcdGetValue(t *testing.T) {
	Convey("With an invalid host, the default value should be returned", t, func() {
		etcd := NewEtcd("http://127.0.0.1:1234")
		value, err := etcd.GetValue("sample_key", "value")

		So(value, ShouldEqual, "")
		So(err, ShouldEqual, EtcdLookupFailure)
	})

	Convey("With a valid host and an unset key, the default value should be returned", t, func() {
		etcd := NewEtcd("http://127.0.0.1:4001")

		value, err := etcd.GetValue("sample_key", "value")
		So(err, ShouldBeNil)
		So(value, ShouldEqual, "value")
	})

	Convey("With a valid host and an set key, the stored value should be returned", t, func() {
		key := "sample_key_" + uuid.New()
		etcd := NewEtcd("http://127.0.0.1:4001")

		err := etcd.SetValue(key, "correct value")
		So(err, ShouldBeNil)

		value, err := etcd.GetValue(key, "value")
		So(value, ShouldEqual, "correct value")
		So(err, ShouldBeNil)
	})
}
