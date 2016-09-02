// Package donitdoc contains document modeling for all user's data and
// brings all models together in one store service
// Every data is stored in MongoDB in terms of documents.
//
// Modeling in subpackages does not care about relationships between data and
// who owns them.
// Though, some model can still add additional linkages, but linkages should be
// mainly maintained by upper level packages.
//
// Also check out github.com/iocat/donit/internal/donitlink for more info
package donitdoc
