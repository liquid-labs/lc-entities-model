# Internal Keys

## Overview

The original data design envisioned using internal (backend-private) and external (public) keys. Specifically, each entity had a unique serial internal ID and a unique public UUID ID. Long term, this design is still of interest, but short term, we're going to simplify.

## Why internal keys?

Performance. UUIDs make larger and slower indexes. Mainly. Internal keys might also make data migration and modification easier.

But, as always, this is also a question of scale. The benefit of faster JOINs (etc.) are offset by the complication of needing to translate between the external and internal keys.

## All UUIDs for now.

The main issue is go-pg does not currently support 'query hooks' (or alternatives) that we would need to cleanly implement our dual-key approach. See [this issue](https://github.com/go-pg/pg/issues/1345) for some details.

Dropping the dual-key regime, we can simplify the data model and use pg-orm as intended and (mostly?) do away with data transfer layer.

## Future

Once we have some scale and can test performance of JOINs under both regimes, we can re-evaluate the benefit. We can also give time to pg-orm to get the necessary features.

## References

* [UUID or GUID as Primary Keys? Be Careful!](https://tomharrisonjr.com/uuid-or-guid-as-primary-keys-be-careful-7b2aa3dcb439) This article seems pretty influential and I believe was a primary driver in the original decision to use the dual-key model. Upon further review and digging deeper, it's not clear that his argument is supported empirically. I.e., it sounds good, but may not be as strong as it seems. The strongest point in practice may be "primary keys get around". What we really want is to understand the impact with large, many table joins.
* [MySQL Insert Performance](http://kccoder.com/mysql/uuid-vs-int-insert-performance/) The original analysis was on MySQL, where the impact of UUID seems to be much more significant.
* [INT4 VS INT8 VS UUID VS NUMERIC PERFORMANCE ON BIGGER JOINS](https://www.cybertec-postgresql.com/en/int4-vs-int8-vs-uuid-vs-numeric-performance-on-bigger-joins/) (on Postgres) the performance impact seems to be smaller. This, and other tests, are all pretty trivial, though. In our system, joining 6+ tables is not uncommon.
* In the lc-entities-model project the last version to use the int/ext key regime was 25d7de8ddefcc4727dd840736f2630bd848d8eae.

## Appendix: Unsent writeup to go-pg maintainer.

### For

#### Greater simplify the backend

One of the main goals of the system (which looks very close) is to do away (mostly) with 'data transfer layer' and use the ORM to go straight from the model DB to the Go model without the need to call intervening functions. With a query hook, this seems immanently possible.

While this particular use case may not be dispositive, the general idea of having Go-level models which are composed from a fully normalized DB seems strong. The best way to compose with SQL is JOIN >> SELECT, and modifying the query implicitly through hooks (or field labels?) to add the proper JOINs would make the ORM far more flexible and work with relational math the way it was meant to work.

In my particular case, the problem outlined is a general problem in that all incoming API requests will only have the public (UUID) IDs because that's all the users see, while all tables internally relate using the backend-private internal (int) IDs.

So, doing a select to populate the internal IDs could mean many selects per-object. Adding 3-4 selects and increasing the number of queries by 400-500% per update/insert seems like a pretty stiff penalty in performance and code complexity.

#### JOIN vs SELECT performance

A big motivation for the current setup was to avoid JOINs on primary keys. The original DB was designed for MySQL where the impact of UUID vs INT JOINs _appears_ to be more significant. Performance impact on Postgres _may_ not be as significant. There's plenty of tests with simple table-to-table joins that show relatively moderate effects. In this system, which is fully normalized and builds in authorizations, a 6+ table JOIN would be relatively common.

#### Support internal vs external keys

There's lots of opinion on the web that the "real" primary key should be private. This is partly for security, but the more convincing arguments are operational. With unpublished keys, it's generally easier to merge and migrate.

### Against

Having said all that, if the ideal isn't available, it might not be worth trying to work around it. I think I might just switch over to using UUIDs as the only keys and JOIN on them directly. Migrating to the current internal/external model could always be done later.
