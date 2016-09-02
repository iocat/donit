// Package donitdoc contains document modeling for all user's data
// Every data is stored in MongoDB in terms of documents.
//
// Modeling does not care about relationships and linking between data
// Current modeling has no concern about who owns the data.
//
// Though, some model can still add additional linkages, but linkages should be
// mainly maintained by upper level packages.
package donitdoc
