/*
   Truct, Pretty minimal workflow runner.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <ntsk@airmail.cc>
*/

package Cli

import (
	Types "github.com/neotesk/truct/src/types"
);

var Setup = Types.Command {
    Name: "setup",
    ShortDesc: "Runs the 'setup' workflow",
    LongDesc: "This command runs a workflow, It will run the 'setup' workflow.",
    Descriptor: "[...args]",
    Action: func ( args []string, flags map[ string ] bool, adjectives map[ string ] string ) {
        Do.Action( append( []string{ "setup" }, args... ), flags, adjectives );
    },
};

var Run = Types.Command {
    Name: "run",
    ShortDesc: "Runs the 'run' workflow",
    LongDesc: "This command runs a workflow, It will run the 'run' workflow.",
    Descriptor: "[...args]",
    Action: func ( args []string, flags map[ string ] bool, adjectives map[ string ] string ) {
        Do.Action( append( []string{ "run" }, args... ), flags, adjectives );
    },
};