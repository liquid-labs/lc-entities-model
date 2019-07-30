package coremodel

type Container struct {
  Entity
  Contents []*Entity
}

func NewContainer(
    name string,
    description string,
    ownerPubID PublicID,
    publiclyReadable bool) *Container {
  return &Container{*NewEntity(name, description, ownerPubID, publiclyReadable), []*Entity{}}
}

func (c *Container) Clone() *Container {
  return &Container{*c.Entity.Clone(), c.Contents}
}

func (c *Container) CloneNew() *Container {
  return &Container{*c.Entity.CloneNew(), c.Contents}
}
