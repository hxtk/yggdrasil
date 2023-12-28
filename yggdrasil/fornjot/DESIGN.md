# Fornjót (Draft)

Named for one of the oldest giants (Jotünn) of Norse Mythology, this is an attempt at implementing Google's Colossus.


## Executive Summary

## In Scope

- Implement the GFS API primitives necessary to support RocksDB/LevelDB/Pebble.
- Support concurrent append to data files
- Support file versioning
- Support durability controls
- Support ACL evaluation without SpiceDB (SpiceDB will be hosted on this file system).
- Support for migrating to a higher level of indirection: if a cluster grows dangerously close to its capacity, it must be possible to migrate its data onto a new Fornjót instance that uses the existing one for metadata storage, or if that is not possible then it must be possible to create a new Fornjót instance that has a higher capacity and merge it with the existing cluster without downtime.
## Out of Scope

## Background

Like Fornjót, this component will be the progenitor of a number of other services.

At Google, Colossus is used to host BigTable (which is used to host Colossus, recursively, to fan out to the necessary size). Here, we will use [RocksDB](https://rocksdb.org/docs/getting-started.html) and [Pebble](https://github.com/cockroachdb/pebble), two compatible, embeddable file-based persistent key-value stores, as our substitute for BigTable. Both of those projects use LSMTrees, and will need to be forked and modified to use the client library for this application as their backing file store.

Similarly, Cloud Spanner is built on Colossus as well, and Zanzibar built on top of Spanner.

Likewise, this file system will need to be usable as a NAS to serve the needs of less sophisticated clients, such as backups.

## High-Level Design

[Etcd](https://etcd.io) will be used as the root metadata store for the "root" instance of Fornjót, similar to how Chubby is used as the metadata store for the root instance of Colossus.

The system will be built on top of itself, recursively, until the capacity reaches the desired size. 

Due to the coupling that exists between Colossus and BigTable, a reimplementation of Colossus requires a reimplementation of [BigTable](https://research.google/pubs/pub27898/). Pebble will be used as the SSTable reader for the Tablet store.

Fornjót must, therefore, be able to read its metadata from either etcd or BigTable metadata stores.
## Detailed Design

Much of the design can be found in Google papers and blog posts:

https://research.google/pubs/pub51/
https://cloud.google.com/blog/products/storage-data-transfer/a-peek-behind-colossus-googles-file-system
https://github.com/CodeBear801/tech_summary/blob/master/tech-summary/papers/colossus.md

## Special Constraints
