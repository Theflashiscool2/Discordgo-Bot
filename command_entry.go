package main

// commandEntry is a command entry used to generate the help command.
type commandEntry [3]string

// Name returns the name of the command.
func (c commandEntry) Name() string {
	return c[0]
}

// Category returns the category of the command.
func (c commandEntry) Category() string {
	return c[1]
}

// Description returns the description of the command.
func (c commandEntry) Description() string {
	return c[2]
}
