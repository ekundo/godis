package client

import "github.com/ekundo/godis/resp"

type Command struct {
	// command name
	Name string

	// command arity specification
	//
	// Command arity follows a simple pattern:
	// - positive if command has fixed number of required arguments.
	// - negative if command has minimum number of required arguments, but may have more.
	// Command arity includes counting the command name itself.
	Arity int

	// command flags
	Flags []string

	// position of first key in argument list
	FirstKeyIndex int

	// position of last key in argument list
	LastKeyIndex int

	// step count for locating repeating keys
	KeyStep int
}

// Command returns array of supported command details.
func (c *Client) Command() ([]Command, error) {
	req := cmd([]string{"command"})
	res, err := c.processRequest(req)
	if err != nil {
		return nil, err
	}

	arr, ok := res.(*resp.Array)
	if !ok {
		return nil, UnexpectedResponseError{}
	}
	cmds := arr.Items
	ret := make([]Command, 0, len(cmds))
	for _, cmd := range cmds {
		arr, ok = cmd.(*resp.Array)
		if !ok {
			return nil, UnexpectedResponseError{}
		}
		fields := arr.Items
		if len(fields) != 6 {
			return nil, UnexpectedResponseError{}
		}
		name, err := c.simpleString(fields[0])
		if err != nil {
			return nil, err
		}
		arity, err := c.integer(fields[1])
		if err != nil {
			return nil, err
		}

		arr, ok = fields[2].(*resp.Array)
		if !ok {
			return nil, UnexpectedResponseError{}
		}
		cmdFlags := arr.Items
		flags := make([]string, 0, len(cmdFlags))
		for _, cmdFlag := range cmdFlags {
			flag, err := c.simpleString(cmdFlag)
			if err != nil {
				return nil, err
			}
			flags = append(flags, flag)
		}

		firstKeyIndex, err := c.integer(fields[3])
		if err != nil {
			return nil, err
		}

		lastKeyIndex, err := c.integer(fields[4])
		if err != nil {
			return nil, err
		}

		keyStep, err := c.integer(fields[5])
		if err != nil {
			return nil, err
		}
		ret = append(ret, Command{
			Name:          name,
			Arity:         arity,
			Flags:         flags,
			FirstKeyIndex: firstKeyIndex,
			LastKeyIndex:  lastKeyIndex,
			KeyStep:       keyStep})
	}
	return ret, nil
}
