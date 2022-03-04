# go-queue-sql (github.com/antonio-alexander/go-queue-sql)

The goal of this library is a proof of concept that provides another implementation of [go-queue](github.com/antonio-alexander/go-queue-sql) that is neither in-memory (like finite/infinite) nor on-disk like [go-durostore](github.com/antonio-alexander/go-durostore), but uses a database as a back-end.

I am FULLY aware that the database is NOT the best solution for this problem, but...a database scales considerably better than file systems or in-memory.
