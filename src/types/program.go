/*
    Truct, Pretty minimal workflow runner.
    Open-Source, WTFPL License.

    Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Types

// This is for having a record of arguments and flags
type Argument struct {
    Name string;
    ShortDesc string;
    DefaultValue string;
}

type Flag struct {
    Name string;
    ShortDesc string;
    DefaultValue bool;
}

// This is for storing default arguments
type DefaultArgs struct {
    Arguments []Argument;
    Flags []Flag;
}

// This is for exporting arguments
type ArgumentsList struct {
    Arguments map[ string ] string;
    Flags map[ string ] bool;
    Keywords []string;
}

// This is for registering commands and program itself
type Command struct {
    Name string;
    ShortDesc string;
    LongDesc string;
    Descriptor string;
    Action func( args []string, flags map[ string ] bool, adjectives map[ string ] string );
}

type Program struct {
    Name string;
    Desc string;
    Footer string;
    Commands []Command
    DefaultArgs DefaultArgs
}