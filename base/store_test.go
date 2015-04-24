package base

import (
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2/bson"
)

func (s *BaseSuite) TestStore_Insert(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	err := st.Insert(p)
	c.Assert(err, IsNil)
	c.Assert(p.IsNew(), Equals, false)

	r, err := st.Find(NewQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_InsertOld(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	err := st.Insert(p)
	c.Assert(err, IsNil)

	err = st.Insert(p)
	c.Assert(err, Equals, NonNewDocumentErr)
}

func (s *BaseSuite) TestStore_Update(c *C) {
	p := &Person{FirstName: "foo"}

	st := NewStore(s.conn, "test")
	st.Insert(p)
	st.Insert(&Person{FirstName: "bar"})

	p.FirstName = "qux"
	err := st.Update(p)
	c.Assert(err, IsNil)

	q := NewQuery()
	q.AddCriteria("firstname", "qux")

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
}

func (s *BaseSuite) TestStore_UpdateNew(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")

	err := st.Update(p)
	c.Assert(err, Equals, NewDocumentErr)
}

func (s *BaseSuite) TestStore_Delete(c *C) {
	p := &Person{FirstName: "foo"}
	st := NewStore(s.conn, "test")
	st.Insert(p)

	err := st.Delete(p)
	c.Assert(err, IsNil)

	r, err := st.Find(NewQuery())
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 0)
}

func (s *BaseSuite) TestStore_FindLimit(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewQuery()
	q.Limit = 1
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_FindSkip(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewQuery()
	q.Skip = 1
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "bar")
}

func (s *BaseSuite) TestStore_FindSort(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewQuery()
	q.Sort = Sort{{IdField, Desc}}
	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 2)
	c.Assert(result[0].FirstName, Equals, "bar")
	c.Assert(result[1].FirstName, Equals, "foo")
}

func (s *BaseSuite) TestStore_RawUpdate(c *C) {
	st := NewStore(s.conn, "test")
	st.Insert(&Person{FirstName: "foo"})
	st.Insert(&Person{FirstName: "bar"})

	q := NewQuery()
	q.AddCriteria("firstname", "foo")

	err := st.RawUpdate(q, bson.M{"firstname": "qux"})
	c.Assert(err, IsNil)

	q = NewQuery()
	q.AddCriteria("firstname", "qux")

	r, err := st.Find(q)
	c.Assert(err, IsNil)

	var result []*Person
	c.Assert(r.All(&result), IsNil)
	c.Assert(result, HasLen, 1)
	c.Assert(result[0].FirstName, Equals, "qux")
}
