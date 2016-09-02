// Package donitdoc contains document modeling for all user's data
// Every data is stored in MongoDB in terms of documents.
//
// Modeling does not care about relationships between data and who owns them.
// Though, some model can still add additional linkages, but linkages should be
// mainly maintained by upper level packages.
//
// Also check github.com/iocat/donit/internal/donitlink for more info
package donitdoc
