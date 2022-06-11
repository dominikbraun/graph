package graph

// Hash is a hashing function that takes a vertex of type T and returns a hash value of type K.
//
// Every graph has one particular hashing function and uses that function to retrieve the hash
// values of its vertices. You can either use one of the predefined hashing functions, or, if you
// want to store a custom data type, provide your own function:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// The types T and K used by the hashing function also define the types T and K of Graph.
type Hash[K comparable, T any] func(T) K

// StringHash is a hashing function that accepts a string and returns that string as a hash value
// at the same time.
func StringHash(v string) string {
	return v
}

// IntHash is a hashing function that accepts a int and returns that int as a hash value at the
// same time.
func IntHash(v int) int {
	return v
}

// Int32Hash is a hashing function that accepts a int32 and returns that int32 as a hash value at
// the same time.
func Int32Hash(v int32) int32 {
	return v
}

// Int64Hash is a hashing function that accepts a int64 and returns that int64 as a hash value at
// the same time.
func Int64Hash(v int64) int64 {
	return v
}

// Uint32Hash is a hashing function that accepts a uint32 and returns that uint32 as a hash value
// at the same time.
func Uint32Hash(v uint32) uint32 {
	return v
}

// Uint64Hash is a hashing function that accepts a uint64 and returns that uint64 as a hash value
// at the same time.
func Uint64Hash(v uint64) uint64 {
	return v
}
