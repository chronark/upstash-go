package upstash

type KV struct {
	Key   string
	Value string
}

// The SET command supports a set of options that modify its behavior:
type SetOptions struct {

	// Set the specified expire time, in seconds.
	EX int

	//  Set the specified expire time, in milliseconds.
	PX int

	//  Set the specified Unix time at which the key will expire, in seconds.
	EXAT int

	//  Set the specified Unix time at which the key will expire, in milliseconds.
	PXAT int

	//  Only set the key if it does not already exist.
	NX bool

	//  Only set the key if it already exist.
	XX bool

	//  Retain the time to live associated with the key.
	KEEPTTL bool

	//  Return the old string stored at key, or nil if key did not exist. An error
	//  is returned and SET aborted if the value stored at key is not a string.
	GET bool
}

// The GETEX command supports a set of options that modify its behavior
// Only one of these should be set.
type GetEXOptions struct {
	// Set the specified expire time, in seconds.
	EX int

	//  Set the specified expire time, in milliseconds.
	PX int

	//  Set the specified Unix time at which the key will expire, in seconds.
	EXAT int

	//  Set the specified Unix time at which the key will expire, in milliseconds.
	PXAT int

	//  Remove the time to live associated with the key.
	PERSIST bool
}
