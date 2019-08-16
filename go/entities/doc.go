// This entities package covers basic model and database functionnality. Entities should generally not be directly created/retrieved/etc.
//
// Notice that the methods use the interface 'orm.DB', which accepts either a pg.DB or pg.Tx. This will typically be a Tx (as entity-row changes should be coordinated with other row changes in a transaction).
package entities
