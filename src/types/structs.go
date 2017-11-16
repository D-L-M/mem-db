package types


// Document indices need to store both the document JSON byte array and an
// inverted index of the keys where its entries in the inverted search index
// can be found
type DocumentIndex struct {
	Document []byte
	InvertedKeys []string
}