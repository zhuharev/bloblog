package bloblog

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestAll(t *testing.T) {
	os.Remove("test.bl")
	Convey("test all", t, func() {
		var (
			bl *BlobLog
			e  error
			id int64
		)

		Convey("test create", func() {
			bl, e = New("test.bl")
			So(e, ShouldBeNil)
		})

		Convey("test insert", func() {
			bl, e = New("test.bl")
			id, e = bl.Prepare(2)
			So(e, ShouldBeNil)
			So(id, ShouldEqual, 1)

			e = bl.Write(id, []byte{1, 2})
			So(e, ShouldBeNil)

			e = bl.Write(id, []byte{1, 2, 3})
			So(e, ShouldNotBeNil)
		})

		Convey("test get", func() {
			bl, e = New("test.bl")
			bts, e := bl.Get(1)
			So(e, ShouldBeNil)
			So(len(bts), ShouldEqual, 2)
			So(bts[0], ShouldEqual, 1)
			So(bts[1], ShouldEqual, 2)
		})

	})
}
