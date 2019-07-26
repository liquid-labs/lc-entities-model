package coremodel

type Container struct {
  Entity
}

func NewContainer(
    name string,
    description string,
    ownerPubID PublicID,
    publiclyReadable bool) *Container {
  return &Container{*NewEntity(name, description, ownerPubID, publiclyReadable)}
}

func (c *Container) Clone() *Container {
  return &Container{*c.Entity.Clone()}
}

func (c *Container) CloneNew() *Container {
  return &Container{*c.Entity.CloneNew()}
}
