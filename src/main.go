/*
   Truct, Pretty minimal workflow manager.
   Open-Source, WTFPL License.

   Copyright (C) 2025-20xx Neo <neotesk@airmail.cc>
*/

package main

import (
	Cli "github.com/neotesk/truct/src/cli"
	Types "github.com/neotesk/truct/src/types"
	Internal "github.com/neotesk/truct/src/util"
)

func main () {
    // Setup default argument parameters
    defaultArgs := Types.DefaultArgs {
        Flags: []Types.Flag {
            {
                Name: "d",
                ShortDesc: "Enables Debug Mode",
                DefaultValue: false,
            },
            {
                Name: "s",
                ShortDesc: "Enables Silent Mode",
                DefaultValue: false,
            },
        },
        Arguments: []Types.Argument {
            {
                Name: "filename",
                ShortDesc: "Override workflow file name",
                DefaultValue: "truct.yaml",
            },
        },
    };

    // After setting the default arguments, let's
    // feed them inside the Arguments function
    // so we can get a good output of what we have
    // in our hands.
    argsList := Arguments( defaultArgs );

    // After gathering the output, we will
    // create a program
    program := Types.Program {
        Name: "truct",
        Desc: "Truct is a minimal workflow manager, giving you the ability to store and run your custom workflows for projects.",
        Footer: "For more information, please visit https://github.com/neotesk/truct/wiki",
        DefaultArgs: defaultArgs,
        Commands: Cli.Commands,
    };

    Internal.IsDebug = argsList.Flags[ "d" ];

    // Now we run the program.
    RunProgram( program, argsList );
}